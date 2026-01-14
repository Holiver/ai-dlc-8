import apiService from './apiService';

export interface Product {
  id: number;
  name: string;
  image_url: string;
  points_required: number;
  stock_quantity: number;
  status: string;
  created_at: string;
  updated_at: string;
}

const productService = {
  // Get all active products
  getProducts: async (): Promise<Product[]> => {
    const response = await apiService.get<{ products: Product[] }>('/products');
    return response.data.products;
  },

  // Get product by ID
  getProductById: async (id: number): Promise<Product> => {
    const response = await apiService.get<{ product: Product }>(`/products/${id}`);
    return response.data.product;
  },
};

export default productService;
