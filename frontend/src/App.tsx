import { useEffect, lazy, Suspense } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { useAuthStore } from "./store/authStore";
import { LoadingSpinner } from "./components/common/LoadingSpinner";

// React Query configuration
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      gcTime: 1000 * 60 * 30, // 30 minutes (formerly cacheTime)
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

// Lazy load pages for code splitting
const Login = lazy(() => import("./pages/Login").then(m => ({ default: m.Login })));
const Register = lazy(() => import("./pages/Register").then(m => ({ default: m.Register })));
const Home = lazy(() => import("./pages/Home").then(m => ({ default: m.Home })));
const ProductDetail = lazy(() => import("./pages/ProductDetail").then(m => ({ default: m.ProductDetail })));
const PurchaseProduct = lazy(() => import("./pages/PurchaseProduct").then(m => ({ default: m.PurchaseProduct })));
const Purchases = lazy(() => import("./pages/Purchases").then(m => ({ default: m.Purchases })));
const ShippingLabel = lazy(() => import("./pages/ShippingLabel").then(m => ({ default: m.ShippingLabel })));
const AutoPurchaseWatches = lazy(() => import("./pages/AutoPurchaseWatches").then(m => ({ default: m.AutoPurchaseWatches })));
const Messages = lazy(() => import("./pages/Messages").then(m => ({ default: m.Messages })));
const Chat = lazy(() => import("./pages/Chat").then(m => ({ default: m.Chat })));
const CreateProduct = lazy(() => import("./pages/CreateProduct").then(m => ({ default: m.CreateProduct })));
const AICreateProduct = lazy(() => import("./pages/AICreateProduct"));
const Profile = lazy(() => import("./pages/Profile").then(m => ({ default: m.Profile })));
const Leaderboard = lazy(() => import("./pages/Leaderboard").then(m => ({ default: m.Leaderboard })));
const Favorites = lazy(() => import("./pages/Favorites").then(m => ({ default: m.Favorites })));
const Notifications = lazy(() => import("./pages/Notifications").then(m => ({ default: m.Notifications })));
const AdminDashboard = lazy(() => import("./pages/AdminDashboard").then(m => ({ default: m.AdminDashboard })));

// Protected route wrapper
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated, token, user } = useAuthStore();

  // Check if we have authentication data
  const localToken = localStorage.getItem("auth_token");
  const hasValidAuth = (isAuthenticated && token && user) || localToken;

  if (!hasValidAuth) {
    console.log("Auth check failed:", {
      isAuthenticated,
      hasToken: !!token,
      hasUser: !!user,
      hasLocalToken: !!localToken,
    });
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

function App() {
  const { isAuthenticated, token, user } = useAuthStore();

  // Debug: log auth state
  useEffect(() => {
    console.log("Auth state:", {
      isAuthenticated,
      hasToken: !!token,
      hasUser: !!user,
    });
    const localToken = localStorage.getItem("auth_token");
    console.log("LocalStorage token:", !!localToken);
  }, [isAuthenticated, token, user]);

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Toaster position="top-right" />
        <Suspense
          fallback={
            <div className="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-dark">
              <LoadingSpinner type="spinner" size="xl" text="読み込み中..." />
            </div>
          }
        >
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />

        {/* Protected routes */}
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <Home />
            </ProtectedRoute>
          }
        />
        <Route
          path="/products/:id"
          element={
            <ProtectedRoute>
              <ProductDetail />
            </ProtectedRoute>
          }
        />
        <Route
          path="/purchase/:id"
          element={
            <ProtectedRoute>
              <PurchaseProduct />
            </ProtectedRoute>
          }
        />
        <Route
          path="/purchases"
          element={
            <ProtectedRoute>
              <Purchases />
            </ProtectedRoute>
          }
        />
        <Route
          path="/purchases/:purchaseId/shipping-label"
          element={
            <ProtectedRoute>
              <ShippingLabel />
            </ProtectedRoute>
          }
        />
        <Route
          path="/auto-purchases"
          element={
            <ProtectedRoute>
              <AutoPurchaseWatches />
            </ProtectedRoute>
          }
        />
        <Route
          path="/messages"
          element={
            <ProtectedRoute>
              <Messages />
            </ProtectedRoute>
          }
        />
        <Route
          path="/chat/:id"
          element={
            <ProtectedRoute>
              <Chat />
            </ProtectedRoute>
          }
        />
        <Route
          path="/create"
          element={
            <ProtectedRoute>
              <CreateProduct />
            </ProtectedRoute>
          }
        />
        <Route
          path="/ai/create"
          element={
            <ProtectedRoute>
              <AICreateProduct />
            </ProtectedRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          }
        />
        <Route
          path="/leaderboard"
          element={
            <ProtectedRoute>
              <Leaderboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/favorites"
          element={
            <ProtectedRoute>
              <Favorites />
            </ProtectedRoute>
          }
        />
        <Route
          path="/notifications"
          element={
            <ProtectedRoute>
              <Notifications />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin"
          element={
            <ProtectedRoute>
              <AdminDashboard />
            </ProtectedRoute>
          }
        />

            {/* Redirect to home by default */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Suspense>
      </BrowserRouter>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}

export default App;
