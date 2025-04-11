import { ImagePlus } from "lucide-react";
import { useRef, useState } from "react";

export default function ImageUploader() {
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) {
      return;
    }

    const file = files[0];
    if (file && file.type.startsWith('image/')) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreviewUrl(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    fileInputRef.current?.click();
  };

  return (
    <div
      className="flex flex-col items-center"
    >
      <label
        className="text-black font-medium mb-2"
        htmlFor="picture"
      >
        picture
      </label>
      <input
        id="picture"
        name="picture"
        type="file"
        ref={fileInputRef}
        accept="image/*"
        className="hidden"
        onChange={handleFileSelect}
        onClick={(e) => e.stopPropagation()}
      />
      <div
        className="w-30 h-30 border-2 border-dashed border-gray-300 rounded-full flex items-center justify-center cursor-pointer overflow-hidden"
        onClick={handleClick}
      >
        {previewUrl ? (
          <img
            src={previewUrl}
            alt="preview"
            className="max-w-full max-h-full object-contain"
          />
        ) : (
          <ImagePlus
            className="w-10 h-10 text-gray-300"
          />
        )}
      </div>
    </div>
  );
}