import { NavLink } from "react-router";
import Sidebar from "./ui/Sidebar";

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

export default function HubSidebar() {
  return (
    <Sidebar>
      <ul>
        <li>
          <NavLink className="w-full p-2 bg-gray-500" to="/">
          </NavLink>
        </li>
        <li>
          <NavLink className="w-full p-2 bg-gray-500" to="/">
          </NavLink>
        </li>
        <li>
          <NavLink className="w-full p-2 bg-gray-500" to="/">
          </NavLink>
        </li>
      </ul>
    </Sidebar>
  );
}