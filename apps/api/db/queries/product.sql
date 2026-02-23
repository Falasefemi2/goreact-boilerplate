-- name: CreateProduct :one
INSERT INTO products (user_id, name, description, price, stock)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListProductsByUser :many
SELECT * FROM products
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET
    name        = COALESCE($1, name),
    description = COALESCE($2, description),
    price       = COALESCE($3, price),
    stock       = COALESCE($4, stock),
    updated_at  = NOW()
WHERE id = $5 AND user_id = $6
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1 AND user_id = $2;
