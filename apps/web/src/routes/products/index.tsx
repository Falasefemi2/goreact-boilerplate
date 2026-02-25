import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useProducts,
  useCreateProduct,
  useDeleteProduct,
  useUpdateProduct,
} from "../../hooks/useProducts";
import { useLogout } from "../../hooks/useAuth";
import {
  Trash2,
  Package,
  Plus,
  LogOut,
  AlertCircle,
  Loader,
  Pencil,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogClose,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogPanel,
  DialogPopup,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Field, FieldLabel } from "@/components/ui/field";

import type { ApiError, Product } from "../../types/api";

export const Route = createFileRoute("/products/")({
  component: ProductsPage,
});

function ProductsPage() {
  const navigate = useNavigate();
  const { data, isLoading, error } = useProducts();
  const createProduct = useCreateProduct();
  const deleteProduct = useDeleteProduct();
  const updateProduct = useUpdateProduct();
  const logout = useLogout();

  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [stock, setStock] = useState(0);
  const [description, setDescription] = useState("");
  const [formError, setFormError] = useState("");
  const [open, setOpen] = useState(false);

  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const [editId, setEditId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [editPrice, setEditPrice] = useState("");
  const [editStock, setEditStock] = useState(0);
  const [editDescription, setEditDescription] = useState("");
  const [editFormError, setEditFormError] = useState("");
  const [editOpen, setEditOpen] = useState(false);

  const products = data?.data ?? [];

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    setFormError("");

    if (!name || !price) {
      setFormError("Name and price are required");
      return;
    }

    try {
      await createProduct.mutateAsync({ name, price, stock, description });
      setName("");
      setPrice("");
      setStock(0);
      setDescription("");
      setOpen(false);
    } catch (err) {
      const apiError = err as ApiError;
      setFormError(apiError.error ?? "Failed to create product");
    }
  }

  function openDeleteDialog(id: string) {
    setDeleteId(id);
    setDeleteOpen(true);
  }

  function openEditDialog(product: Product) {
    setEditId(product.id);
    setEditName(product.name);
    setEditPrice(product.price);
    setEditStock(product.stock);
    setEditDescription(
      product.description.Valid ? product.description.String : "",
    );
    setEditFormError("");
    setEditOpen(true);
  }

  async function handleEdit(e: React.FormEvent) {
    e.preventDefault();
    setEditFormError("");

    if (!editName || !editPrice || !editId) {
      setEditFormError("Name and price are required");
      return;
    }

    try {
      await updateProduct.mutateAsync({
        id: editId,
        data: {
          name: editName,
          description: editDescription,
          price: editPrice,
          stock: editStock,
        },
      });
      setEditOpen(false);
      setEditId(null);
    } catch (err) {
      const apiError = err as ApiError;
      setEditFormError(apiError.error ?? "Failed to update product");
    }
  }

  async function handleDelete() {
    if (!deleteId) return;

    try {
      await deleteProduct.mutateAsync(deleteId);
      setDeleteOpen(false);
      setDeleteId(null);
    } catch {
      // optional error handling
    }
  }

  async function handleLogout() {
    await logout.mutateAsync();
    navigate({ to: "/login" });
  }

  return (
    <div className="max-w-4xl mx-auto py-10 space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Package className="w-6 h-6" />
          <h1 className="text-2xl font-semibold">Products</h1>
        </div>

        <Button variant="destructive" onClick={handleLogout}>
          <LogOut className="w-4 h-4 mr-2" />
          Logout
        </Button>
      </div>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger render={<Button />}>
          <Plus className="w-4 h-4 mr-2" />
          Add Product
        </DialogTrigger>

        <DialogPopup className="sm:max-w-md">
          <Form onSubmit={handleCreate} className="contents">
            <DialogHeader>
              <DialogTitle>Create Product</DialogTitle>
              <DialogDescription>
                Fill in the details below to add a new product.
              </DialogDescription>
            </DialogHeader>

            <DialogPanel className="grid gap-4">
              <Field>
                <FieldLabel>Product Name</FieldLabel>
                <Input value={name} onChange={(e) => setName(e.target.value)} />
              </Field>

              <Field>
                <FieldLabel>Product Description</FieldLabel>
                <Input
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </Field>

              <Field>
                <FieldLabel>Price</FieldLabel>
                <Input
                  value={price}
                  onChange={(e) => setPrice(e.target.value)}
                />
              </Field>

              <Field>
                <FieldLabel>Stock</FieldLabel>
                <Input
                  type="number"
                  value={stock}
                  onChange={(e) => setStock(Number(e.target.value))}
                />
              </Field>

              {formError && (
                <Alert>
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{formError}</AlertDescription>
                </Alert>
              )}
            </DialogPanel>

            <DialogFooter>
              <DialogClose render={<Button variant="ghost" />}>
                Cancel
              </DialogClose>

              <Button type="submit" disabled={createProduct.isPending}>
                {createProduct.isPending ? (
                  <>
                    <Loader className="w-4 h-4 mr-2 animate-spin" />
                    Creating...
                  </>
                ) : (
                  "Create Product"
                )}
              </Button>
            </DialogFooter>
          </Form>
        </DialogPopup>
      </Dialog>

      <Card>
        <CardHeader>
          <CardTitle>Product List</CardTitle>
          <CardDescription>Manage your available products</CardDescription>
        </CardHeader>

        <CardContent>
          {isLoading ? (
            <div className="flex items-center gap-2">
              <Loader className="animate-spin w-4 h-4" />
              <span>Loading products...</span>
            </div>
          ) : error ? (
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>Failed to load products</AlertDescription>
            </Alert>
          ) : products.length === 0 ? (
            <p className="text-muted-foreground">
              No products yet. Add one above.
            </p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Price</TableHead>
                  <TableHead>Stock</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>

              <TableBody>
                {products.map((product) => (
                  <TableRow key={product.id}>
                    <TableCell>{product.name}</TableCell>
                    <TableCell>
                      {product.description.Valid
                        ? product.description.String
                        : "—"}
                    </TableCell>
                    <TableCell>
                      <Badge>${product.price}</Badge>
                    </TableCell>
                    <TableCell>{product.stock}</TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          size="icon"
                          variant="ghost"
                          onClick={() => openEditDialog(product)}
                        >
                          <Pencil className="w-4 h-4" />
                        </Button>
                        <Button
                          size="icon"
                          variant="destructive"
                          onClick={() => openDeleteDialog(product.id)}
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <Dialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <DialogPopup className="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>Delete Product</DialogTitle>
            <DialogDescription>
              This action cannot be undone. This will permanently delete the
              selected product.
            </DialogDescription>
          </DialogHeader>

          <DialogPanel>
            <p className="text-sm text-muted-foreground">
              Are you sure you want to continue?
            </p>
          </DialogPanel>

          <DialogFooter>
            <DialogClose render={<Button variant="ghost" />}>
              Cancel
            </DialogClose>

            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteProduct.isPending}
            >
              {deleteProduct.isPending ? (
                <>
                  <Loader className="w-4 h-4 mr-2 animate-spin" />
                  Deleting...
                </>
              ) : (
                "Delete"
              )}
            </Button>
          </DialogFooter>
        </DialogPopup>
      </Dialog>

      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogPopup className="sm:max-w-md">
          <Form onSubmit={handleEdit} className="contents">
            <DialogHeader>
              <DialogTitle>Edit Product</DialogTitle>
              <DialogDescription>
                Update the fields below and save your changes.
              </DialogDescription>
            </DialogHeader>

            <DialogPanel className="grid gap-4">
              <Field>
                <FieldLabel>Product Name</FieldLabel>
                <Input
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                />
              </Field>

              <Field>
                <FieldLabel>Product Description</FieldLabel>
                <Input
                  value={editDescription}
                  onChange={(e) => setEditDescription(e.target.value)}
                />
              </Field>

              <Field>
                <FieldLabel>Price</FieldLabel>
                <Input
                  value={editPrice}
                  onChange={(e) => setEditPrice(e.target.value)}
                />
              </Field>

              <Field>
                <FieldLabel>Stock</FieldLabel>
                <Input
                  type="number"
                  value={editStock}
                  onChange={(e) => setEditStock(Number(e.target.value))}
                />
              </Field>

              {editFormError && (
                <Alert>
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{editFormError}</AlertDescription>
                </Alert>
              )}
            </DialogPanel>

            <DialogFooter>
              <DialogClose render={<Button variant="ghost" />}>
                Cancel
              </DialogClose>

              <Button type="submit" disabled={updateProduct.isPending}>
                {updateProduct.isPending ? (
                  <>
                    <Loader className="w-4 h-4 mr-2 animate-spin" />
                    Saving...
                  </>
                ) : (
                  "Save Changes"
                )}
              </Button>
            </DialogFooter>
          </Form>
        </DialogPopup>
      </Dialog>
    </div>
  );
}
