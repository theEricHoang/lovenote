import { Outlet } from "react-router";
import HubSidebar from "~/components/HubSidebar";

export default function HubLayout() {
  return (
    <>
      <HubSidebar />
      <Outlet />
    </>
  );
}