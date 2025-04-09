import { CircleFadingPlus, ImagePlus } from "lucide-react";
import { useState } from "react";
import Modal from "./ui/Modal";
import Button from "./ui/Button";
import ImageUploader from "./ui/ImageUploader";

export default function NewRelationshipDialog() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  
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
        <form>
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
              type="text"
              placeholder="lovey doveys..."
              className="w-full mt-1 px-3 py-2 text-black placeholder:text-gray-300 focus:border-rose-400 focus:outline-none focus:ring-4 focus:ring-rose-50 shadow-sm border border-slate-300 rounded-md"
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
              className="inline"
              type="submit"
            >
              create!
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}