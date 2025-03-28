import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "lovenote" },
    { name: "description", content: "welcome to lovenote" },
  ];
}

export default function Home() {
  return <div className="text-black">hello</div>;
}
