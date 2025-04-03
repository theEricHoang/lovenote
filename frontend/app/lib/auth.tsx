import { createContext, useContext, useEffect, useState } from "react";
import { api, setUpRequestInterceptors, setUpResponseInterceptors } from "./http";

interface User {
  id: number;
  username: string;
  email?: string;
  profilePicture: string;
  bio?: string;
  createdAt?: string;
}

interface AuthContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  accessToken: string | null;
  setAccessToken: (token: string | null) => void;
  isLoading: boolean;
  setIsLoading: (isLoading: boolean) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({
  children
}: {
  children: React.ReactNode
}) {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // expose access token so it can be used outside of React
  const getAccessToken = () => accessToken;

  useEffect(() => {
    const refreshAuth = async () => {
      try {
        const response = await api.post("users/refresh");
        const newAccessToken = response.data.access;
        setAccessToken(newAccessToken);

        const userResponse = await api.get("users/me", {
          headers: {
            Authorization: `Bearer ${newAccessToken}`
          },
        });
        setUser({
          id: userResponse.data.id,
          username: userResponse.data.username,
          email: userResponse.data.email,
          profilePicture: userResponse.data.profile_picture,
          bio: userResponse.data.bio,
        });
      } catch (error) {
        console.error("refresh failed,", error);
        setAccessToken(null);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };

    setUpRequestInterceptors(getAccessToken);
    setUpResponseInterceptors(setAccessToken);

    refreshAuth();
  }, []);

  return (
    <AuthContext.Provider value={{ user, setUser, accessToken, setAccessToken, isLoading, setIsLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return context;
}