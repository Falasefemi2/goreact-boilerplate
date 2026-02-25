import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useLogin } from "../hooks/useAuth";
import type { ApiError } from "../types/api";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardHeader,
  CardPanel,
  CardTitle,
} from "@/components/ui/card";
import { Field, FieldLabel } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  const navigate = useNavigate();
  const login = useLogin();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (!email || !password) {
      setError("Email and password are required");
      return;
    }

    try {
      await login.mutateAsync({ email, password });
      navigate({ to: "/products" });
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.error ?? "Login failed");
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-full max-w-xs">
        <CardHeader>
          <CardTitle>Login to your account</CardTitle>
          <CardAction>
            <Link
              className="text-muted-foreground text-sm hover:underline"
              to="/register"
            >
              Register
            </Link>
          </CardAction>
        </CardHeader>

        <CardPanel>
          <Form onSubmit={handleSubmit}>
            <Field>
              <FieldLabel>Email</FieldLabel>
              <Input
                placeholder="Enter your email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </Field>

            <Field>
              <FieldLabel>Password</FieldLabel>
              <Input
                placeholder="Enter your password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </Field>

            {error && <p className="text-sm text-red-500">{error}</p>}

            <Button className="w-full" type="submit" disabled={login.isPending}>
              {login.isPending ? "Logging in..." : "Login"}
            </Button>
          </Form>
        </CardPanel>
      </Card>
    </div>
  );
}
