import { CircleFadingPlus } from "lucide-react";
import { useRef, useState } from "react";
import Modal from "./ui/Modal";
import Button from "./ui/Button";
import ImageUploader from "./ui/ImageUploader";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "~/lib/http";
import axios from "axios";

export default function NewRelationshipDialog() {
  const queryClient = useQueryClient();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const formRef = useRef<HTMLFormElement>(null);

  const createRelationshipMutation = useMutation({
    mutationFn: async ({ name, file }: { name: string; file: File }) => {
      if (!file || file.size === 0) {
        const createRes = await api.post("/relationships", {
          "name": name,
          "picture": "",
        });

        return createRes.data;
      }

      // get presigned URL
      const presignRes = await api.post("/relationships/presign-put", {
        "filename": file.name,
        "content_type": file.type,
      });
      const { url, key } = await presignRes.data;

      // upload to s3
      await axios.put(url, file, {
        headers: {
          "Content-Type": file.type,
        },
      });
      
      // upload public url and relationship name to backend
      const pictureUrl = `https://lovenote-images.s3.amazonaws.com/${key}`;
      const createRes = await api.post("/relationships", {
        "name": name,
        "picture": pictureUrl,
      });

      return createRes.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['relationships'] })
      closeModal();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!formRef.current) return;

    const formData = new FormData(formRef.current);

    const name = formData.get("name") as string;
    const file = formData.get("picture") as File;

    createRelationshipMutation.mutate({ name, file });
  }

  const openModal = () => setIsModalOpen(true);
  const closeModal = () => setIsModalOpen(false);

  return (
    <>
      <button
        className="flex w-full p-2 rounded-md hover:bg-gray-50 text-black/25 font-medium items-center"
        onClick={openModal}
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

      <Modal
        isOpen={isModalOpen}
        onClose={closeModal}
        title="create a new relationship"
      >
        <form ref={formRef} onSubmit={handleSubmit}>
          <ImageUploader />

          <div className="mt-6">
            <label
              className="text-black font-medium"
              htmlFor="name"
            >
              name
            </label>
            <input
              id="name"
              name="name"
              type="text"
              placeholder="lovey doveys..."
              className="w-full mt-1 px-3 py-2 text-black placeholder:text-gray-300 focus:border-rose-400 focus:outline-none focus:ring-4 focus:ring-rose-50 shadow-sm border border-slate-300 rounded-md"
              required
            />
          </div>

          <div className="flex justify-end mt-6 space-x-2">
            <Button
              className="inline"
              variant="outline"
              onClick={closeModal}
            >
              cancel
            </Button>
            <Button
              type="submit"
              isLoading={createRelationshipMutation.isPending}
            >
              create!
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}