import { useQuery } from "@tanstack/react-query";
import type { Route } from "./+types";
import { api } from "~/lib/http";
import type { User } from "~/lib/auth";
import DragCanvas from "~/components/ui/DragCanvas";
import { Plus } from "lucide-react";

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
      const response = await api.get(`relationships/${params.relationshipId}/notes`);
      return response.data;
    },
  });

  return (
    <>
      <button
        className="p-2 border text-black border-gray-300 bg-white hover:bg-gray-200 fixed top-4 right-4 z-50 rounded-full"
        aria-controls="sidebar"
      >
        <Plus size={20} />
        <span className="sr-only">open new note</span>
      </button>
      <DragCanvas>
        {
          noteQuery.data?.map((note) => (
            <div
              key={note.id}
              style={{
                backgroundColor: note.color,
                position: "absolute",
                left: `${note.position_x}px`,
                top: `${note.position_y}px`,
              }}
              className="text-black w-64 min-h-64 shadow-lg"
            >
              <div className="w-full p-2 font-medium bg-black/3">
                {note.title}
              </div>
              <div className="p-2">
                {note.content}
              </div>
            </div>
          ))
        }
      </DragCanvas>
    </>
  );
}