import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { purchaseService } from '@/services/purchases'
import toast from 'react-hot-toast'

// Query keys
export const purchaseKeys = {
  all: ['purchases'] as const,
  lists: () => [...purchaseKeys.all, 'list'] as const,
  list: () => [...purchaseKeys.lists()] as const,
  details: () => [...purchaseKeys.all, 'detail'] as const,
  detail: (id: string) => [...purchaseKeys.details(), id] as const,
}

// Fetch purchases list
export const usePurchases = (role?: 'buyer' | 'seller') => {
  return useQuery({
    queryKey: purchaseKeys.list(),
    queryFn: () => purchaseService.list(role),
    staleTime: 1000 * 60 * 2, // 2 minutes
  })
}

// Fetch single purchase
export const usePurchase = (id: string) => {
  return useQuery({
    queryKey: purchaseKeys.detail(id),
    queryFn: () => purchaseService.getById(id),
    enabled: !!id,
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

// Create purchase mutation
export const useCreatePurchase = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: purchaseService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: purchaseKeys.lists() })
      toast.success('購入を完了しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '購入に失敗しました')
    },
  })
}

// Complete purchase mutation
export const useCompletePurchase = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (purchaseId: string) => purchaseService.complete(purchaseId),
    onSuccess: (_, purchaseId) => {
      queryClient.invalidateQueries({ queryKey: purchaseKeys.detail(purchaseId) })
      queryClient.invalidateQueries({ queryKey: purchaseKeys.lists() })
      toast.success('取引を完了しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '取引の完了に失敗しました')
    },
  })
}

// Generate shipping label mutation
export const useGenerateShippingLabel = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: purchaseService.generateShippingLabel,
    onSuccess: (_, purchaseId) => {
      queryClient.invalidateQueries({ queryKey: purchaseKeys.detail(purchaseId) })
      toast.success('配送ラベルを作成しました')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || '配送ラベルの作成に失敗しました')
    },
  })
}
