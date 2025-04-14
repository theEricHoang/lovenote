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
        layout("routes/hub/_layout.tsx", [
            index("routes/hub/index.tsx"),
            route("relationships/:relationshipId", "routes/hub/relationship.tsx"),
        ]),
    ]),
] satisfies RouteConfig;
