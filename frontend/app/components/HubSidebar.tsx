import { NavLink } from "react-router";
import Sidebar from "./ui/Sidebar";
import { CircleFadingPlus } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "~/lib/http";

const mockData = [
  {
    id: 1,
    name: "Gangy",
    profilePicture: "https://img.freepik.com/free-vector/gradient-heart_78370-478.jpg",
  },
  {
    id: 2,
    name: "lovely bubblez",
    profilePicture: "https://img.freepik.com/free-vector/gradient-heart_78370-478.jpg",
  },
  {
    id: 3,
    name: "Say Gex",
    profilePicture: "https://img.freepik.com/free-vector/gradient-heart_78370-478.jpg",
  },
];

interface Relationship {
  id: number;
  name: string;
  picture: string;
  created_at: string;
}

export default function HubSidebar() {
  const { data } = useQuery({
    queryKey: ['relationships'],
    queryFn: async (): Promise<Array<Relationship>> => {
      const response = await api.get("/relationships");
      return response.data;
    },
  })

  return (
    <Sidebar>
      <button
        className="flex w-full p-2 rounded-md hover:bg-gray-50 text-black font-medium items-center"
      >
        <CircleFadingPlus
          size={40}
        />
        <span
          className="mx-2"
        >
          new relationship
        </span>
      </button>

      <ul className="space-y-1">
        {data?.map((relationship: Relationship) => (
          <li key={relationship.id}>
            <NavLink className="flex p-2 rounded-md hover:bg-gray-50 text-black font-medium items-center" to="/">
              <img
                src={relationship.picture}
                alt={`${relationship.name} icon`}
                className="w-10 h-10 rounded-full"
              />
              <span className="mx-2">
                {relationship.name}
              </span>
            </NavLink>
          </li>
        ))}
      </ul>
    </Sidebar>
  );
}