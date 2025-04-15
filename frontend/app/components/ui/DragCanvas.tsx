import { useRef } from "react";

export default function DragCanvas({
  className,
  children,
}: {
  className?: string;
  children?: React.ReactNode;
}) {
  const canvasRef = useRef<HTMLDivElement>(null);
  const isDragging = useRef(false);

  const handleMouseDown = (e: React.MouseEvent) => {
    isDragging.current = true;
    canvasRef.current!.dataset.startX = e.clientX.toString();
    canvasRef.current!.dataset.startY = e.clientY.toString();
    canvasRef.current!.dataset.scrollLeft = canvasRef.current!.scrollLeft.toString();
    canvasRef.current!.dataset.scrollTop = canvasRef.current!.scrollTop.toString();
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!isDragging.current || !canvasRef.current) return;

    const dx = e.clientX - parseInt(canvasRef.current.dataset.startX!);
    const dy = e.clientY - parseInt(canvasRef.current.dataset.startY!);

    canvasRef.current.scrollLeft = parseInt(canvasRef.current.dataset.scrollLeft!) - dx;
    canvasRef.current.scrollTop = parseInt(canvasRef.current.dataset.scrollTop!) - dy;
  };

  const handleMouseUp = () => {
    isDragging.current = false;
  };

  return (
    <div
      ref={canvasRef}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
      className="relative w-screen h-screen overflow-hidden cursor-grab active:cursor-grabbing"
    >
      <div
        className="absolute"
        style={{ width: "3840px", height: "2160px" }}
      >
        {children}
      </div>
    </div>
  );
}