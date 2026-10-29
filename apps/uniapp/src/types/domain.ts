export type EntityID = string

export interface FabricColor {
  name: string
  hex: string
}

export interface Fabric {
  id: EntityID
  ownerId: EntityID
  name: string
  material: string
  pattern: string
  styleTags: string[]
  colors: FabricColor[]
  coverAssetId?: EntityID
  coverUrl?: string
  pieces: FabricPiece[]
  memo?: string
  version: number
  createdAt: string
  updatedAt: string
}

export interface FabricPiece {
  id: EntityID
  fabricId: EntityID
  label: string
  lengthMm?: number
  widthMm?: number
  quantity: number
  unit: string
  storageLocation?: string
  purchaseChannel?: string
  purchaseDate?: string
  totalPriceCent?: number
  memo?: string
  status: 'in_stock' | 'used_up' | 'archived'
}

export interface Asset {
  id: EntityID
  storageKey: string
  url: string
}

export interface CreateFabricPayload {
  name: string
  material: string
  pattern: string
  styleTags: string[]
  colors: FabricColor[]
  memo: string
  coverAssetId?: EntityID
  piece?: {
    label: string
    lengthMm?: number
    widthMm?: number
    storageLocation: string
    purchaseChannel: string
    purchaseDate?: string
    totalPriceCent?: number
    memo: string
  }
}

export interface Project {
  id: EntityID
  ownerId: EntityID
  name: string
  kind: string
  status: 'planned' | 'making' | 'completed' | 'archived'
  memo?: string
  version: number
  createdAt: string
  updatedAt: string
}

