import { createContext, useContext, useEffect, useState } from "react";
import { setUpRequestInterceptors, setUpResponseInterceptors } from "./http";

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
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({
  children
}: {
  children: React.ReactNode
}) {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);

  // expose access token so it can be used outside of React
  const getAccessToken = () => accessToken;

  useEffect(() => {
    setUpRequestInterceptors(getAccessToken);
    setUpResponseInterceptors(setAccessToken);
  }, []);

  return (
    <AuthContext.Provider value={{ user, setUser, accessToken, setAccessToken }}>
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