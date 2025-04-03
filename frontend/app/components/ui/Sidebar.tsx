import { useState } from "react";
import { Menu } from "lucide-react";
import Logo from "./Logo";

export default function Sidebar({
  className,
  children,
}: {
  className?: string;
  children?: React.ReactNode;
}) {
  const [isOpen, setIsOpen] = useState(true);

  return (
    <>
      <button
        className="p-2 border text-black border-gray-300 bg-white hover:bg-gray-200 fixed top-4 left-4 z-50 rounded-full"
        onClick={() => { setIsOpen(!isOpen) }}
        aria-expanded={isOpen}
        aria-controls="sidebar"
      >
        <Menu size={20} />
        <span className="sr-only">{isOpen ? 'close menu' : 'open menu'}</span>
      </button>
      <aside
          id="sidebar"
          className={`${className} ${
            isOpen ? 'translate-x-0' : '-translate-x-full'
          } fixed top-0 left-0 h-full border-r border-gray-300 bg-white w-64 transition-transform duration-300 ease-in-out z-40 shadow-xl`}
          aria-label="Main navigation"
      >
        <Logo className="ml-15 mt-4.75" />
        <nav
          className="p-4 space-y-8"
        >
          {children}
        </nav>
      </aside>
    </>
  );
}