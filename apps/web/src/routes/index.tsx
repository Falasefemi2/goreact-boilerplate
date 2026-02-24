import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useMe } from "../hooks/useAuth";

export const Route = createFileRoute("/")({
  component: IndexPage,
});

function IndexPage() {
  const { data, isLoading } = useMe();

  if (isLoading) return <p>Loading...</p>;

  // if logged in go to products, otherwise go to login
  if (data) return <Navigate to="/products" />;
  return <Navigate to="/login" />;
}
