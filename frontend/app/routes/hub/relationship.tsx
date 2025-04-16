import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types";
import { api } from "~/lib/http";
import { useAuth, type User } from "~/lib/auth";
import DragCanvas from "~/components/ui/DragCanvas";
import { Plus, X } from "lucide-react";

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
  const { user } = useAuth();
  const queryClient = useQueryClient();

  const noteQuery = useQuery({
    queryKey: [`relationship${params.relationshipId}notes`],
    queryFn: async (): Promise<Array<Note>> => {
      const response = await api.get(`relationships/${params.relationshipId}/notes`);
      return response.data;
    },
  });

  const deleteNoteMutation = useMutation({
    mutationFn: async ({ noteId }: { noteId: number }) => {
      const res = await api.delete(`/relationships/${params.relationshipId}/notes/${noteId}`);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [`relationship${params.relationshipId}notes`] });
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
              className="text-black w-64 min-h-64 shadow-lg flex flex-col"
            >
              <div className="w-full p-2 font-medium bg-black/3 flex flex-row justify-between">
                <span>
                  {note.title}
                </span>
                {note.author.id == user?.id &&
                  <button
                    className="text-black/10 hover:text-black/30"
                    onClick={() => deleteNoteMutation.mutate({ noteId: note.id })}
                  >
                    <X size={20} />
                  </button>
                }
              </div>

              <div className="p-2">
                {note.content}
              </div>

              <div className="w-full p-2 flex flex-row items-center space-x-2 mt-auto">
                <img
                  className="w-5 h-5 rounded-full"
                  src={note.author.profile_picture}
                  alt={`${note.author.username} profile picture`}
                />
                <span className="text-sm font-medium">
                  {note.author.username}
                </span>
              </div>
            </div>
          ))
        }
      </DragCanvas>
    </>
  );
}