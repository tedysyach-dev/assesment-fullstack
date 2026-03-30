import { useProducts, useDeleteProduct } from "../hooks/useProducts";

const ProductList = () => {
  const { data: products, isLoading, isError, error } = useProducts();
  const { mutate: deleteProduct, isPending: isDeleting } = useDeleteProduct();

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="h-20 bg-gray-100 rounded-lg animate-pulse" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <div className="text-red-500 text-sm p-4 bg-red-50 rounded-lg">
        Gagal memuat produk: {(error as Error).message}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {products?.map((product) => (
        <div
          key={product.id}
          className="flex items-center justify-between p-4 bg-white rounded-lg border border-gray-100 shadow-sm"
        >
          <div>
            <p className="font-medium text-gray-800 text-sm">{product.title}</p>
            <p className="text-gray-400 text-xs mt-0.5">{product.category}</p>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-blue-600 font-semibold text-sm">
              ${product.price}
            </span>
            <button
              onClick={() => deleteProduct(product.id)}
              disabled={isDeleting}
              className="text-xs text-red-400 hover:text-red-600 transition disabled:opacity-50"
            >
              Hapus
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default ProductList;
