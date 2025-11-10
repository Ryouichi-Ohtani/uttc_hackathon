import { useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import { useAuthStore } from "./store/authStore";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { Home } from "./pages/Home";
import { ProductDetail } from "./pages/ProductDetail";
import { PurchaseProduct } from "./pages/PurchaseProduct";
import { Purchases } from "./pages/Purchases";
import { Messages } from "./pages/Messages";
import { Chat } from "./pages/Chat";
import { CreateProduct } from "./pages/CreateProduct";
import AICreateProduct from "./pages/AICreateProduct";
import { Profile } from "./pages/Profile";
import { Leaderboard } from "./pages/Leaderboard";
import { Favorites } from "./pages/Favorites";
import { Notifications } from "./pages/Notifications";

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
    <BrowserRouter>
      <Toaster position="top-right" />
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

        {/* Redirect to home by default */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
