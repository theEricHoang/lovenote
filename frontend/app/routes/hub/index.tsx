import { Navigate } from "react-router";
import { useAuth } from "~/lib/auth";

export default function Hub() {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <div className="text-black">loading...</div>;
  }

  if (!user) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />;
  }

  return (
    <div className="text-black">hello {user?.username}</div>
  );
}