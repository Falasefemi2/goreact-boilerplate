import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useProducts,
  useCreateProduct,
  useDeleteProduct,
} from "../../hooks/useProducts";
import { useLogout } from "../../hooks/useAuth";
import type { ApiError } from "../../types/api";

export const Route = createFileRoute("/products/")({
  component: ProductsPage,
});

function ProductsPage() {
  const navigate = useNavigate();
  const { data, isLoading, error } = useProducts();
  const createProduct = useCreateProduct();
  const deleteProduct = useDeleteProduct();
  const logout = useLogout();

  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [stock, setStock] = useState(0);
  const [formError, setFormError] = useState("");

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    setFormError("");
    try {
      await createProduct.mutateAsync({ name, price, stock });
      setName("");
      setPrice("");
      setStock(0);
    } catch (err) {
      const apiError = err as ApiError;
      setFormError(apiError.error ?? "Failed to create product");
    }
  }

  async function handleLogout() {
    await logout.mutateAsync();
    navigate({ to: "/login" });
  }

  if (isLoading) return <p>Loading products...</p>;
  if (error) return <p>Failed to load products</p>;

  const products = data?.data ?? [];

  return (
    <div style={{ maxWidth: 800, margin: "40px auto", padding: 24 }}>
      <div style={{ display: "flex", justifyContent: "space-between" }}>
        <h1>Products</h1>
        <button onClick={handleLogout}>Logout</button>
      </div>

      {/* Create form */}
      <form onSubmit={handleCreate} style={{ marginBottom: 32 }}>
        <h2>Add Product</h2>
        <input
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <input
          placeholder="Price (e.g. 9.99)"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
          required
        />
        <input
          type="number"
          placeholder="Stock"
          value={stock}
          onChange={(e) => setStock(Number(e.target.value))}
        />
        {formError && <p style={{ color: "red" }}>{formError}</p>}
        <button type="submit" disabled={createProduct.isPending}>
          {createProduct.isPending ? "Creating..." : "Add Product"}
        </button>
      </form>

      {/* Product list */}
      {products.length === 0 ? (
        <p>No products yet. Add one above.</p>
      ) : (
        <table style={{ width: "100%", borderCollapse: "collapse" }}>
          <thead>
            <tr>
              <th>Name</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {products.map((product) => (
              <tr key={product.id}>
                <td>{product.name}</td>
                <td>${product.price}</td>
                <td>{product.stock}</td>
                <td>
                  <button
                    onClick={() => deleteProduct.mutate(product.id)}
                    disabled={deleteProduct.isPending}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
