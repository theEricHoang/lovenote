import { NavLink, useNavigate } from "react-router";
import Sidebar from "./ui/Sidebar";
import { ChevronUp, CircleFadingPlus, LogOut, Settings } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "~/lib/http";
import { useAuth } from "~/lib/auth";
import { useState } from "react";
import NewRelationshipDialog from "./NewRelationshipDialog";

interface Relationship {
  id: number;
  name: string;
  picture: string;
  created_at: string;
}

export default function HubSidebar() {
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
  const { user } = useAuth();
  const navigate = useNavigate();

  const { data } = useQuery({
    queryKey: ['relationships'],
    queryFn: async (): Promise<Array<Relationship>> => {
      const response = await api.get("/relationships");
      return response.data;
    },
  });

  const handleLogout = async () => {
    try {
      await api.post("/users/logout");
      navigate("/");
    } catch (error) {
      console.error("logout error:", error);
    }
  };

  return (
    <Sidebar>
      <nav
        className="mt-4 p-2 space-y-1 flex-grow overflow-y-auto"
      >
        <NewRelationshipDialog />

        <ul className="space-y-1">
          {data?.map((relationship: Relationship) => (
            <li key={relationship.id}>
              <NavLink className="flex p-2 rounded-md hover:bg-gray-50 text-black font-semibold items-center" to="/">
                <img
                  src={relationship.picture}
                  alt={`${relationship.name} icon`}
                  className="w-10 h-10 rounded-full object-cover"
                />
                <span className="mx-2">
                  {relationship.name}
                </span>
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>

      <div
        className="relative border-t border-gray-200 p-1"
      >
        <button
          className="rounded-md hover:bg-gray-50 flex w-full p-2 items-center"
          onClick={() => { setIsUserMenuOpen(!isUserMenuOpen) }}
        >
          <img
            src={user?.profilePicture}
            alt="Current user profile picture"
            className="w-10 h-10 rounded-full"
          />
          <span className="text-black ml-2 flex-grow text-left">{user?.username}</span>
          <ChevronUp size={28} className="text-black"/>
        </button>

        {isUserMenuOpen && (
        <div className="absolute bottom-full left-0 right-0 mb-1 bg-white shadow-lg border border-gray-200 rounded-md text-black overflow-hidden">
          <ul className="py-1">
            <li>
              <button
                className="w-full px-4 py-2 text-left flex items-center hover:bg-gray-100"
              >
                <Settings className="mr-2 h-4 w-4" />
                <span>Settings</span>
              </button>
            </li>
            <li>
              <button
                className="w-full px-4 py-2 text-left flex items-center hover:bg-gray-100"
                onClick={() => handleLogout()}
              >
                <LogOut className="mr-2 h-4 w-4" />
                <span>Logout</span>
              </button>
            </li>
          </ul>
        </div>
      )}
      </div>
    </Sidebar>
  );
}