import { useQuery } from "@tanstack/react-query";
import type { Route } from "./+types";
import { api } from "~/lib/http";
import type { User } from "~/lib/auth";

interface Note {
  id: number;
  author: User;
  title: string;
  content: string;
  position_x: number;
  position_y: number;
  color: string;
  created_at: string;
}

export default function Relationship({
  params
}: Route.ComponentProps) {
  const noteQuery = useQuery({
    queryKey: [`relationship${params.relationshipId}notes`],
    queryFn: async (): Promise<Array<Note>> => {
      const response = await api.get(`relationships/${params.relationshipId}/notes?month=3&year=2025`);
      return response.data;
    },
  });

  return (
    <div className="flex w-full h-full">
      {
        noteQuery.data?.map((note) => (
          <div className={`bg-[${note.color}] p-4`}>

          </div>
        ))
      }
    </div>
  );
}