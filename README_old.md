# EcoMate - Next-Gen Sustainable Flea Market

![EcoMate Logo](docs/logo.png)

## 🌍 Overview

EcoMate is a next-generation flea market application that visualizes environmental impact. Every transaction shows how much CO2 is saved by buying used items instead of new ones, gamifying sustainability with levels, achievements, and leaderboards.

## ✨ Key Features

### Core Features
- 🔐 **User Authentication**: JWT-based secure authentication
- 📦 **Product Listing & Purchase**: Full e-commerce flow
- 💬 **Real-time DM**: WebSocket-powered instant messaging
- 🤖 **AI-Powered Assistance**:
  - Auto-generate product descriptions from images
  - Price suggestions based on market data
  - Inappropriate content detection

### Advanced Features
- 🌱 **CO2 Impact Calculation**: Real-time environmental impact visualization
- 🏆 **Gamification**: Levels, achievements, and sustainability scores
- 📊 **Sustainability Dashboard**: Personal impact tracking with charts
- 🎨 **3D Product Viewer**: Three.js 360° product visualization
- ⚡ **CDN Image Optimization**: Fast image delivery
- 🔄 **gRPC Communication**: High-performance AI service integration

## 🏗️ Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   React     │────→│  Go Backend  │────→│  CloudSQL   │
│  (Vercel)   │←────│  (CloudRun)  │←────│ (PostgreSQL)│
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                           │ gRPC
                           ↓
                    ┌──────────────┐
                    │   Python AI  │
                    │   Service    │
                    │  (CloudRun)  │
                    └──────────────┘
                           │
                           ↓
                    ┌──────────────┐
                    │  Gemini API  │
                    └──────────────┘
```

### Tech Stack

**Backend (Go)**
- Framework: Gin
- Architecture: Clean Architecture (Domain/UseCase/Infrastructure)
- Database: PostgreSQL with GORM
- Auth: JWT
- Real-time: Gorilla WebSocket
- Testing: testify

**AI Service (Python)**
- Framework: FastAPI
- AI: Google Gemini API
- Orchestration: LangChain/LangGraph
- Communication: gRPC
- ML: Custom CO2 calculation models

**Frontend (React)**
- Framework: React 18 + TypeScript
- State: Zustand
- Styling: Tailwind CSS
- 3D: Three.js + React Three Fiber
- Charts: Recharts
- HTTP: Axios
- WebSocket: native WebSocket API

**Infrastructure**
- Cloud: Google Cloud Platform
- Backend Deploy: Cloud Run
- Frontend Deploy: Vercel
- Database: Cloud SQL (PostgreSQL)
- CDN: Cloud CDN / Cloud Storage
- CI/CD: GitHub Actions

## 📁 Project Structure

```
ecomate/
├── backend/               # Go backend
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── domain/       # Domain layer (entities)
│   │   ├── usecase/      # Use case layer (business logic)
│   │   ├── infrastructure/ # Infrastructure layer (DB, external APIs)
│   │   ├── interfaces/   # Interface layer (HTTP handlers)
│   │   └── config/
│   ├── migrations/
│   ├── go.mod
│   └── Dockerfile
│
├── ai-service/           # Python AI service
│   ├── app/
│   │   ├── main.py
│   │   ├── services/
│   │   │   ├── gemini.py
│   │   │   ├── co2_calculator.py
│   │   │   └── langchain_workflow.py
│   │   ├── grpc_server/
│   │   └── models/
│   ├── proto/
│   ├── requirements.txt
│   └── Dockerfile
│
├── frontend/             # React frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── common/
│   │   │   ├── auth/
│   │   │   ├── products/
│   │   │   ├── messages/
│   │   │   └── sustainability/
│   │   ├── pages/
│   │   ├── hooks/
│   │   ├── store/
│   │   ├── services/
│   │   ├── types/
│   │   └── App.tsx
│   ├── public/
│   ├── package.json
│   └── vite.config.ts
│
├── proto/               # Shared protobuf definitions
│   └── product_analysis.proto
│
├── docs/                # Documentation
│   ├── database-design.md
│   ├── api-specification.md
│   └── ui-ux-design.md
│
├── .github/
│   └── workflows/
│       ├── backend-ci.yml
│       └── frontend-ci.yml
│
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Python 3.11+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose
- Google Cloud SDK

### Environment Setup

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/ecomate.git
cd ecomate
```

2. **Backend setup**
```bash
cd backend
cp .env.example .env
# Edit .env with your configuration
go mod download
go run cmd/api/main.go
```

3. **AI Service setup**
```bash
cd ai-service
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
python app/main.py
```

4. **Frontend setup**
```bash
cd frontend
npm install
npm run dev
```

### Docker Compose (Recommended for development)

```bash
docker-compose up
```

This will start:
- Backend on `http://localhost:8080`
- AI Service on `http://localhost:8000`
- Frontend on `http://localhost:3000`
- PostgreSQL on `localhost:5432`

## 🧪 Testing

**Backend**
```bash
cd backend
go test ./... -v
```

**AI Service**
```bash
cd ai-service
pytest
```

**Frontend**
```bash
cd frontend
npm test
```

## 📊 Database Migrations

```bash
cd backend
go run cmd/migrate/main.go up
```

## 🌐 Deployment

### Backend (Cloud Run)
```bash
cd backend
gcloud run deploy ecomate-api \
  --source . \
  --region asia-northeast1 \
  --allow-unauthenticated
```

### AI Service (Cloud Run)
```bash
cd ai-service
gcloud run deploy ecomate-ai \
  --source . \
  --region asia-northeast1 \
  --allow-unauthenticated
```

### Frontend (Vercel)
```bash
cd frontend
vercel --prod
```

## 🎯 Demo Scenario

1. **Sign Up**: Create account with email/password
2. **List Item**: Upload product image → AI generates description & suggests price → Review CO2 impact
3. **Browse**: Explore products with sustainability filters
4. **Purchase**: Buy item → See CO2 savings → Earn achievement
5. **Dashboard**: View total impact, level up, check leaderboard
6. **3D View**: Rotate and inspect product in 3D

## 📈 Environmental Impact Calculation

CO2 savings are calculated using:
- Product category baseline emissions
- Manufacturing country (shipping distance)
- Product age (degradation factor)
- Condition multiplier

Formula:
```
CO2_saved = (NEW_PRODUCTION_EMISSIONS + SHIPPING_EMISSIONS) - (USED_SHIPPING_EMISSIONS + DEGRADATION)
```

## 🏆 Achievements

- **First Step**: Complete first transaction
- **Eco Warrior**: Save 10kg CO2
- **Planet Hero**: Save 50kg CO2
- **Climate Champion**: Save 100kg CO2
- **Master Trader**: Complete 50 transactions

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) file

## 👥 Team

- **Backend Lead**: [Name]
- **AI Engineer**: [Name]
- **Frontend Lead**: [Name]
- **DevOps**: [Name]

## 🙏 Acknowledgments

- Google Gemini API for AI capabilities
- Open source community
- UTTC Hackathon organizers

---

Built with 💚 for the planet
