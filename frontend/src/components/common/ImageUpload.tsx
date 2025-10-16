import { useState, useRef } from 'react'
import { Button } from './Button'
import api from '@/services/api'
import toast from 'react-hot-toast'

interface ImageUploadProps {
  onUpload: (urls: string[]) => void
  maxImages?: number
  existingImages?: string[]
}

export const ImageUpload = ({
  onUpload,
  maxImages = 10,
  existingImages = []
}: ImageUploadProps) => {
  const [uploading, setUploading] = useState(false)
  const [images, setImages] = useState<string[]>(existingImages)
  const [previews, setPreviews] = useState<string[]>(existingImages)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || [])

    if (files.length === 0) return

    if (images.length + files.length > maxImages) {
      toast.error(`最大${maxImages}枚まで選択できます`)
      return
    }

    // Validate file types
    const invalidFiles = files.filter(file => !file.type.startsWith('image/'))
    if (invalidFiles.length > 0) {
      toast.error('画像ファイルのみアップロード可能です')
      return
    }

    // Validate file sizes (max 5MB per file)
    const oversizedFiles = files.filter(file => file.size > 5 * 1024 * 1024)
    if (oversizedFiles.length > 0) {
      toast.error('ファイルサイズは5MB以下にしてください')
      return
    }

    // Create preview URLs
    const newPreviews = files.map(file => URL.createObjectURL(file))
    setPreviews([...previews, ...newPreviews])

    // Upload files
    setUploading(true)
    try {
      const formData = new FormData()
      files.forEach(file => {
        formData.append('images', file)
      })

      const response = await api.post('/v1/upload/images', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      const uploadedUrls = response.data.images.map((img: any) => img.url)
      const newImages = [...images, ...uploadedUrls]
      setImages(newImages)
      onUpload(newImages)

      toast.success(`${files.length}枚の画像をアップロードしました`)
    } catch (error: any) {
      console.error('Upload error:', error)
      toast.error('画像のアップロードに失敗しました')
      // Remove previews on error
      setPreviews(previews)
    } finally {
      setUploading(false)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
  }

  const handleRemoveImage = (index: number) => {
    const newImages = images.filter((_, i) => i !== index)
    const newPreviews = previews.filter((_, i) => i !== index)
    setImages(newImages)
    setPreviews(newPreviews)
    onUpload(newImages)
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4">
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          multiple
          onChange={handleFileSelect}
          className="hidden"
        />
        <Button
          type="button"
          variant="outline"
          onClick={() => fileInputRef.current?.click()}
          disabled={uploading || images.length >= maxImages}
        >
          {uploading ? '📤 アップロード中...' : '📷 画像を選択'}
        </Button>
        <span className="text-sm text-gray-600">
          {images.length} / {maxImages} 枚
        </span>
      </div>

      {previews.length > 0 && (
        <div className="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
          {previews.map((preview, index) => (
            <div key={index} className="relative group">
              <img
                src={preview}
                alt={`Preview ${index + 1}`}
                className="w-full h-32 object-cover rounded-lg border border-gray-200"
              />
              <button
                type="button"
                onClick={() => handleRemoveImage(index)}
                className="absolute top-2 right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center opacity-0 group-hover:opacity-100 transition"
              >
                ×
              </button>
              {index === 0 && (
                <div className="absolute bottom-2 left-2 bg-primary-500 text-white text-xs px-2 py-1 rounded">
                  メイン
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      <p className="text-xs text-gray-500">
        ※ 最初の画像がメイン画像として表示されます。JPG, PNG, GIF形式に対応（各5MB以下）
      </p>
    </div>
  )
}
