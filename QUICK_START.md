# EcoMate - Quick Start Guide 🚀

## Prerequisites

- Docker & Docker Compose
- Node.js 18+ (for local frontend development)
- Go 1.21+ (for local backend development)
- Python 3.11+ (for local AI service development)
- Google Gemini API Key

## 🏃‍♂️ Quick Start (5 minutes)

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd UTTC_hackathon
```

### 2. Environment Variables

Create `.env` file in the root:

```bash
# Create .env file
cat > .env << EOF
GOOGLE_API_KEY=your-gemini-api-key-here
EOF
```

### 3. Start All Services

```bash
# Start everything with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f
```

### 4. Initialize Database

```bash
# The backend will auto-migrate tables on startup
# Check if migration succeeded
docker-compose logs backend | grep "Database connected"
```

### 5. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Health Check**: http://localhost:8080/health

## 📱 Test the Application

### 1. Register a User

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@ecomate.com",
    "username": "eco_demo",
    "password": "password123",
    "display_name": "Eco Demo User"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@ecomate.com",
    "password": "password123"
  }'

# Save the token from response
export TOKEN="<your-jwt-token>"
```

### 3. Create a Product (with AI assistance)

```bash
# Using multipart form data
curl -X POST http://localhost:8080/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -F "title=Vintage Uniqlo Fleece Jacket" \
  -F "description=" \
  -F "price=800" \
  -F "category=clothing" \
  -F "condition=good" \
  -F "use_ai_assistance=true" \
  -F "images=@./sample-jacket.jpg"
```

### 4. List Products

```bash
curl http://localhost:8080/v1/products
```

## 🔧 Local Development

### Backend (Go)

```bash
cd backend

# Install dependencies
go mod download

# Setup environment
cp .env.example .env
# Edit .env with your configuration

# Run locally
go run cmd/api/main.go
```

### AI Service (Python)

```bash
cd ai-service

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Setup environment
cp .env.example .env
# Add your GOOGLE_API_KEY

# Generate proto files
python -m grpc_tools.protoc \
  -I../proto \
  --python_out=./proto \
  --grpc_python_out=./proto \
  ../proto/product_analysis.proto

# Run gRPC server
python -m app.grpc_server.server
```

### Frontend (React)

```bash
cd frontend

# Install dependencies
npm install

# Setup environment
echo "VITE_API_URL=http://localhost:8080/v1" > .env.local

# Run development server
npm run dev

# Build for production
npm run build
```

## 🧪 Testing

### Backend Tests

```bash
cd backend
go test ./... -v
```

### Frontend Tests

```bash
cd frontend
npm test
```

## 📊 Database Management

### Access PostgreSQL

```bash
# Connect to database
docker exec -it ecomate-db psql -U ecomate -d ecomate_db

# List tables
\dt

# View users
SELECT id, email, username, sustainability_score FROM users;

# Exit
\q
```

### Reset Database

```bash
docker-compose down -v
docker-compose up -d
```

## 🌐 Deployment

### Cloud Run (Backend & AI Service)

```bash
# Login to GCP
gcloud auth login

# Set project
gcloud config set project YOUR_PROJECT_ID

# Deploy Backend
cd backend
gcloud run deploy ecomate-api \
  --source . \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --set-env-vars="DB_HOST=<cloud-sql-ip>,JWT_SECRET=<secret>,AI_SERVICE_URL=<ai-service-url>"

# Deploy AI Service
cd ../ai-service
gcloud run deploy ecomate-ai \
  --source . \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --set-env-vars="GOOGLE_API_KEY=<your-key>"
```

### Vercel (Frontend)

```bash
cd frontend

# Install Vercel CLI
npm install -g vercel

# Deploy
vercel --prod

# Set environment variable
vercel env add VITE_API_URL
# Enter: https://ecomate-api-xxx.run.app/v1
```

### Cloud SQL Setup

```bash
# Create instance
gcloud sql instances create ecomate-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=asia-northeast1

# Create database
gcloud sql databases create ecomate_db --instance=ecomate-db

# Create user
gcloud sql users create ecomate \
  --instance=ecomate-db \
  --password=<secure-password>

# Get connection name
gcloud sql instances describe ecomate-db --format="value(connectionName)"
```

## 🐛 Troubleshooting

### Backend not connecting to database

```bash
# Check database is running
docker ps | grep postgres

# Check connection from backend
docker exec -it ecomate-backend ping postgres

# View backend logs
docker-compose logs backend
```

### AI Service errors

```bash
# Check API key is set
docker exec -it ecomate-ai env | grep GOOGLE_API_KEY

# Test AI service directly
docker exec -it ecomate-ai python -c "import os; print(os.getenv('GOOGLE_API_KEY'))"

# View AI service logs
docker-compose logs ai-service
```

### Frontend not loading

```bash
# Check if backend is accessible
curl http://localhost:8080/health

# Check frontend logs
docker-compose logs frontend

# Rebuild frontend
docker-compose up -d --build frontend
```

### gRPC connection issues

```bash
# Test gRPC connection from backend to AI service
docker exec -it ecomate-backend ping ai-service

# Check AI service is listening
docker exec -it ecomate-ai netstat -tulpn | grep 50051
```

## 📈 Performance Optimization

### Database Indexing

All critical indexes are already defined in the schema:
- User email/username lookups
- Product search (GIN index for full-text)
- Conversation/message queries
- Sustainability logs by user

### CDN Setup (Production)

```bash
# Create Cloud Storage bucket
gsutil mb -l asia-northeast1 gs://ecomate-products

# Set CORS
gsutil cors set cors-config.json gs://ecomate-products

# Setup Cloud CDN
gcloud compute backend-buckets create ecomate-cdn \
  --gcs-bucket-name=ecomate-products \
  --enable-cdn
```

### Caching Strategy

- Static assets: 1 year cache (immutable)
- API responses: Use React Query with stale-while-revalidate
- Images: CDN with edge caching

## 🎯 Demo Checklist

Before Demo Day, ensure:

- [ ] Database seeded with demo data
- [ ] All services deployed to cloud
- [ ] Test user account created
- [ ] Sample products uploaded with AI-generated descriptions
- [ ] WebSocket messaging working
- [ ] Dashboard showing CO2 stats
- [ ] 3D viewer functional (if model available)
- [ ] Mobile responsive design tested
- [ ] Presentation slides prepared
- [ ] Demo flow practiced

## 🔗 Useful Commands

```bash
# View all container logs
docker-compose logs -f

# Restart specific service
docker-compose restart backend

# Rebuild and restart
docker-compose up -d --build

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Check service health
curl http://localhost:8080/health
curl http://localhost:3000/health
```

## 📚 Additional Resources

- [Full Implementation Guide](IMPLEMENTATION_GUIDE.md)
- [Database Schema](docs/database-design.md)
- [API Documentation](docs/api-specification.md)
- [UI/UX Design](docs/ui-ux-design.md)

## 🆘 Support

If you encounter issues:

1. Check the logs: `docker-compose logs -f`
2. Verify environment variables are set
3. Ensure all ports are available (3000, 8080, 5432, 50051)
4. Try rebuilding: `docker-compose up -d --build`
5. Reset everything: `docker-compose down -v && docker-compose up -d`

---

Happy coding! 🌱 Let's build a sustainable future together!
