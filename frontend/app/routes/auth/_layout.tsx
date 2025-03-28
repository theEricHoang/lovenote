import { Outlet } from "react-router";

export default function Layout() {
  return (
    <div className="flex justify-between">
      <div>layout</div>
      <Outlet />
    </div>
  );
}