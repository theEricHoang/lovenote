import { type RouteConfig, index, layout, prefix, route } from "@react-router/dev/routes";

export default [
    layout("routes/layout.tsx", [
        index("routes/index.tsx"),
    ]),
    layout("routes/auth/_layout.tsx", [
        route("/login", "routes/auth/login.tsx"),
        route("/register", "routes/auth/register.tsx"),
    ]),
    ...prefix("hub", [
        index("routes/hub/index.tsx"),
    ]),
] satisfies RouteConfig;
