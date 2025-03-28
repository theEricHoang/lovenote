import { Outlet } from "react-router";
import NavBar from "~/components/ui/Navbar";

export default function Layout() {
  return (
  <>
    <NavBar />
    <Outlet />
  </>
  );
}