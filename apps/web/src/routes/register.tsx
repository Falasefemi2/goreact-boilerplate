import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useRegister } from "../hooks/useAuth";
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

export const Route = createFileRoute("/register")({
  component: RegisterPage,
});

function RegisterPage() {
  const navigate = useNavigate();
  const register = useRegister();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      await register.mutateAsync({ email, password });
      navigate({ to: "/products" });
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.error ?? "Registration failed");
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-full max-w-xs">
        <CardHeader>
          <CardTitle>Register Your account</CardTitle>
          <CardAction>
            <Link
              to="/login"
              className="text-muted-foreground text-sm hover:underline"
            >
              Login
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

            <Button
              className="w-full"
              type="submit"
              disabled={register.isPending}
            >
              {register.isPending ? "Creating account..." : "Register"}
            </Button>
          </Form>
        </CardPanel>
      </Card>
    </div>
  );
}
