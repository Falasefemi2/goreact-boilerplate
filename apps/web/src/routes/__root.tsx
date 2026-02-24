import { createRootRoute, Outlet, Navigate } from "@tanstack/react-router";
import { useMe } from "../hooks/useAuth";

export const Route = createRootRoute({
  component: RootLayout,
});

function RootLayout() {
  return (
    <div>
      <Outlet />
    </div>
  );
}
