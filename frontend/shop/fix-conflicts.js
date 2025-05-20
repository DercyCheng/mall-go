#!/usr/bin/env node

import fs from 'fs';
import path from 'path';

const fixes = [
    {
        file: '/Users/dercyc/go/src/pro/mall-go/frontend/shop/src/views/checkout/index.vue',
        replacements: [
            {
                from: 'import { generateConfirmOrder, submitOrder } from \'@/api/order\';',
                to: 'import { generateConfirmOrder, submitOrder as submitOrderApi } from \'@/api/order\';'
            },
            {
                from: 'const submitOrder = async () => {',
                to: 'const handleSubmitOrder = async () => {'
            },
            {
                from: '@click="submitOrder"',
                to: '@click="handleSubmitOrder"'
            },
            {
                from: 'const result = await submitOrder(orderData);',
                to: 'const result = await submitOrderApi(orderData);'
            }
        ]
    },
    {
        file: '/Users/dercyc/go/src/pro/mall-go/frontend/shop/src/views/product/detail.vue',
        replacements: [
            {
                from: 'import { addToCart } from \'@/api/cart\';',
                to: 'import { addToCart as addToCartApi } from \'@/api/cart\';'
            },
            {
                from: 'const addToCart = async () => {',
                to: 'const handleAddToCart = async () => {'
            },
            {
                from: '@click="addToCart"',
                to: '@click="handleAddToCart"'
            },
            {
                from: 'await addToCart({',
                to: 'await addToCartApi({'
            }
        ]
    },
    {
        file: '/Users/dercyc/go/src/pro/mall-go/frontend/shop/src/views/user/favorites.vue',
        replacements: [
            {
                from: 'import { getFavoritesList, removeFromFavorites } from \'@/api/user\';',
                to: 'import { getFavoritesList, removeFromFavorites as removeFromFavoritesApi } from \'@/api/user\';'
            },
            {
                from: 'import { addToCart } from \'@/api/cart\';',
                to: 'import { addToCart as addToCartApi } from \'@/api/cart\';'
            },
            {
                from: 'const addToCart = async (productId) => {',
                to: 'const handleAddToCart = async (productId) => {'
            },
            {
                from: 'const removeFromFavorites = async (productId) => {',
                to: 'const handleRemoveFromFavorites = async (productId) => {'
            },
            {
                from: '@click="addToCart(item.id)"',
                to: '@click="handleAddToCart(item.id)"'
            },
            {
                from: '@click="removeFromFavorites(item.id)"',
                to: '@click="handleRemoveFromFavorites(item.id)"'
            },
            {
                from: 'await addToCart({',
                to: 'await addToCartApi({'
            },
            {
                from: 'await removeFromFavorites(productId);',
                to: 'await removeFromFavoritesApi(productId);'
            }
        ]
    },
    {
        file: '/Users/dercyc/go/src/pro/mall-go/frontend/shop/src/views/user/order-detail.vue',
        replacements: [
            {
                from: 'import { getOrderDetail, cancelOrder, deleteOrder, confirmReceiveOrder } from \'@/api/order\';',
                to: 'import { getOrderDetail, cancelOrder as cancelOrderApi, deleteOrder as deleteOrderApi, confirmReceiveOrder } from \'@/api/order\';'
            },
            {
                from: 'const cancelOrder = async () => {',
                to: 'const handleCancelOrder = async () => {'
            },
            {
                from: 'const deleteOrder = async () => {',
                to: 'const handleDeleteOrder = async () => {'
            },
            {
                from: '@click="cancelOrder"',
                to: '@click="handleCancelOrder"'
            },
            {
                from: '@click="deleteOrder"',
                to: '@click="handleDeleteOrder"'
            },
            {
                from: 'await cancelOrder(order.value.id);',
                to: 'await cancelOrderApi(order.value.id);'
            },
            {
                from: 'await deleteOrder(order.value.id);',
                to: 'await deleteOrderApi(order.value.id);'
            }
        ]
    },
    {
        file: '/Users/dercyc/go/src/pro/mall-go/frontend/shop/src/views/user/profile.vue',
        replacements: [
            {
                from: 'import { getUserInfo, updateUserInfo, updatePassword } from \'@/api/user\';',
                to: 'import { getUserInfo, updateUserInfo, updatePassword as updatePasswordApi } from \'@/api/user\';'
            },
            {
                from: 'const updatePassword = async () => {',
                to: 'const handleUpdatePassword = async () => {'
            },
            {
                from: '@click="updatePassword"',
                to: '@click="handleUpdatePassword"'
            },
            {
                from: 'await updatePassword(',
                to: 'await updatePasswordApi('
            }
        ]
    }
];

// Apply all the fixes
fixes.forEach((fix) => {
    try {
        // Read the file
        const filePath = fix.file;
        let content = fs.readFileSync(filePath, 'utf8');

        // Apply all replacements for this file
        fix.replacements.forEach((replacement) => {
            content = content.replace(replacement.from, replacement.to);
        });

        // Write the modified content back to the file
        fs.writeFileSync(filePath, content);
        console.log(`Fixed: ${filePath}`);
    } catch (error) {
        console.error(`Error fixing ${fix.file}:`, error.message);
    }
});

console.log('All fixes applied!');
