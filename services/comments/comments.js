import request from '../../utils/request';

/**
 * 获取商品评论
 * @param {number} productId 商品ID
 * @param {number} page 页码
 * @param {number} limit 每页数量
 */
export function fetchProductComments(productId, page = 1, limit = 10) {
    return request.get(`/api/v1/comments/product/${productId}`, {
        page,
        limit
    }).then(response => {
        return response.list.map(comment => ({
            id: comment.id,
            userName: comment.user?.nick_name || '用户',
            avatar: comment.user?.avatar || '',
            rating: comment.rating,
            content: comment.content,
            images: comment.images ? comment.images.split(',') : [],
            createTime: comment.created_at
        }));
    }).catch(error => {
        console.error('获取商品评论失败:', error);
        return [];
    });
}

/**
 * 创建商品评论
 * @param {object} commentData 评论数据
 */
export function createComment(commentData) {
    return request.post('/api/v1/comments', {
        product_id: commentData.productId,
        order_id: commentData.orderId,
        rating: commentData.rating,
        content: commentData.content,
        images: commentData.images ? commentData.images.join(',') : ''
    });
}
