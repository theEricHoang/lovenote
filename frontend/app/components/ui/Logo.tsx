import { Link } from "react-router";

export default function Logo({
  className,
}: {
  className?: string;
}) {
  return (
    <Link to="/" className={`${className} flex items-center space-x-2`}>
      <svg 
        xmlns="http://www.w3.org/2000/svg" 
        fill="red" 
        viewBox="0 0 24 24" 
        strokeWidth="1.5" 
        stroke="currentColor" 
        className="w-8 h-8"
      >
        <path 
          strokeLinecap="round" 
          strokeLinejoin="round" 
          d="M12 21C12 21 4 13.73 4 8.35 4 5.42 6.42 3 9.35 3c1.54 0 3.03.76 3.91 1.94C14.62 3.76 16.11 3 17.65 3 20.58 3 23 5.42 23 8.35 23 13.73 15 21 12 21z"
        />
      </svg>
      <span className="text-2xl font-bold text-gray-900">lovenote</span>
    </Link>
  );
}