import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { adminService } from '@/services/admin'
import toast from 'react-hot-toast'

// Query keys
export const adminKeys = {
  all: ['admin'] as const,
  users: () => [...adminKeys.all, 'users'] as const,
  products: () => [...adminKeys.all, 'products'] as const,
  purchases: () => [...adminKeys.all, 'purchases'] as const,
}

// Fetch users
export const useAdminUsers = () => {
  return useQuery({
    queryKey: adminKeys.users(),
    queryFn: () => adminService.getUsers(),
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

// Fetch products
export const useAdminProducts = () => {
  return useQuery({
    queryKey: adminKeys.products(),
    queryFn: () => adminService.getProducts(),
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

// Fetch purchases
export const useAdminPurchases = () => {
  return useQuery({
    queryKey: adminKeys.purchases(),
    queryFn: () => adminService.getPurchases(),
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

// Update user mutation
export const useUpdateUser = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ userId, data }: { userId: string; data: any }) =>
      adminService.updateUser(userId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.users() })
      toast.success('ユーザー情報を更新しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '更新に失敗しました')
    },
  })
}

// Delete user mutation
export const useDeleteUser = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: adminService.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.users() })
      toast.success('ユーザーを削除しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || 'ユーザーの削除に失敗しました')
    },
  })
}

// Update product mutation
export const useAdminUpdateProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ productId, data }: { productId: string; data: any }) =>
      adminService.updateProduct(productId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.products() })
      toast.success('商品情報を更新しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '更新に失敗しました')
    },
  })
}

// Delete product mutation
export const useAdminDeleteProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: adminService.deleteProduct,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.products() })
      toast.success('商品を削除しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '商品の削除に失敗しました')
    },
  })
}

// Update purchase mutation
export const useAdminUpdatePurchase = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ purchaseId, data }: { purchaseId: string; data: any }) =>
      adminService.updatePurchase(purchaseId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.purchases() })
      toast.success('購入情報を更新しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '更新に失敗しました')
    },
  })
}
