import { Outlet } from "react-router";
import Logo from "~/components/ui/Logo";

export default function Layout() {
  return (
    <div className="flex size-full items-start justify-center sm:h-screen">
      <div className="hidden size-full min-h-min items-center justify-center bg-rose-50 p-12 lg:flex">
        <Logo />
      </div>
      <Outlet />
    </div>
  );
}