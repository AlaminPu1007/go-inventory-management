
-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, product_id, quantity, price
    )
VALUES (
    $1, $2, $3, $4
    )
RETURNING *;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, product_id, quantity, price, status
FROM order_items
WHERE order_id = $1
ORDER BY id;


-- name: UpdateOrderItemStatus :exec
UPDATE order_items
SET status = $2
WHERE id = $1;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

-- name: ListOrderItemsByUser :many
SELECT 
    oi.id AS order_item_id,
    oi.order_id,
    oi.product_id,
    oi.quantity,
    oi.price,
    oi.created_at AS item_created_at,
    oi.updated_at AS item_updated_at,

    o.status,
    o.total_amount,
    o.created_at AS order_created_at,
    o.updated_at AS order_updated_at
FROM order_items oi
INNER JOIN orders o ON oi.order_id = o.id
WHERE o.user_id = $1
ORDER BY oi.id
LIMIT $2 OFFSET $3;
