import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { productService } from '@/services/products'
import { Product, ProductFilters } from '@/types'
import toast from 'react-hot-toast'

// Query keys
export const productKeys = {
  all: ['products'] as const,
  lists: () => [...productKeys.all, 'list'] as const,
  list: (filters: ProductFilters) => [...productKeys.lists(), filters] as const,
  details: () => [...productKeys.all, 'detail'] as const,
  detail: (id: string) => [...productKeys.details(), id] as const,
  favorites: () => [...productKeys.all, 'favorites'] as const,
}

// Fetch products list
export const useProducts = (filters: ProductFilters) => {
  return useQuery({
    queryKey: productKeys.list(filters),
    queryFn: async () => {
      try {
        const data = await productService.list(filters)
        return data.products
      } catch (error: any) {
        toast.error('商品の読み込みに失敗しました')
        throw error
      }
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

// Fetch single product
export const useProduct = (id: string) => {
  return useQuery({
    queryKey: productKeys.detail(id),
    queryFn: () => productService.getById(id),
    enabled: !!id,
    staleTime: 1000 * 60 * 10, // 10 minutes
  })
}

// Fetch user favorites (using product list with filter)
export const useFavorites = () => {
  return useQuery({
    queryKey: productKeys.favorites(),
    queryFn: async () => {
      // Fetch products and filter favorites client-side
      const data = await productService.list({ limit: 100 })
      return data.products.filter(product => product.is_favorited)
    },
    staleTime: 1000 * 60 * 3, // 3 minutes
  })
}

// Create product mutation
export const useCreateProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
      toast.success('商品を出品しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '出品に失敗しました')
    },
  })
}

// Update product mutation
export const useUpdateProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) =>
      productService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: productKeys.detail(variables.id) })
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
      toast.success('商品を更新しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '更新に失敗しました')
    },
  })
}

// Delete product mutation
export const useDeleteProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
      toast.success('商品を削除しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '削除に失敗しました')
    },
  })
}

// Toggle favorite mutation
export const useToggleFavorite = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ productId, isFavorited }: { productId: string; isFavorited: boolean }) => {
      if (isFavorited) {
        await productService.removeFavorite(productId)
      } else {
        await productService.addFavorite(productId)
      }
      return { productId, isFavorited: !isFavorited }
    },
    onMutate: async ({ productId, isFavorited }) => {
      // Optimistic update
      await queryClient.cancelQueries({ queryKey: productKeys.detail(productId) })

      const previousProduct = queryClient.getQueryData<Product>(productKeys.detail(productId))

      if (previousProduct) {
        queryClient.setQueryData<Product>(productKeys.detail(productId), {
          ...previousProduct,
          is_favorited: !isFavorited,
        })
      }

      return { previousProduct }
    },
    onError: (_error: any, variables, context) => {
      // Rollback on error
      if (context?.previousProduct) {
        queryClient.setQueryData(productKeys.detail(variables.productId), context.previousProduct)
      }
      toast.error('お気に入りの更新に失敗しました')
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: productKeys.favorites() })
      toast.success(data.isFavorited ? 'お気に入りに追加しました' : 'お気に入りから削除しました')
    },
  })
}
