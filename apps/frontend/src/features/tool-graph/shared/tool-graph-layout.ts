import type { ToolGraphResponseOutput } from '@/api/generated/model/ToolGraphResponse.zod'

import type { ToolGraphLayoutMode } from './tool-graph-types'

export type GraphPoint = {
  x: number
  y: number
}

const CATEGORY_ROW: Record<ToolGraphResponseOutput['nodes'][number]['category'], number> = {
  language: 0,
  framework: 1,
  library: 2,
}

export const kindLabel = (value: ToolGraphResponseOutput['links'][number]['kind']): string =>
  value.replaceAll('_', ' ')

export const kindColor = (value: ToolGraphResponseOutput['links'][number]['kind']): string => {
  switch (value) {
    case 'built_on':
      return '#4fb8b2'
    case 'inspired_by':
      return '#6ec89a'
    case 'alternative_to':
      return '#f59e0b'
    case 'replaced_by':
      return '#f97316'
    case 'used_with':
      return '#3b82f6'
    default: {
      const neverValue: never = value

      return neverValue
    }
  }
}

export const computeNodePositions = (
  graph: ToolGraphResponseOutput,
  layoutMode: ToolGraphLayoutMode,
): Map<string, GraphPoint> => {
  if (layoutMode === 'chronological') {
    const sorted = [...graph.nodes].sort((a, b) => a.releaseYear - b.releaseYear)
    const minYear = sorted[0]?.releaseYear ?? new Date().getFullYear()

    return new Map(
      sorted.map((node) => [
        node.id,
        {
          x: (node.releaseYear - minYear) * 190 + 120,
          y: CATEGORY_ROW[node.category] * 170 + 110,
        },
      ]),
    )
  }

  const focusIndex = graph.nodes.findIndex((node) => node.id === graph.focusNodeId)
  const focusNode = graph.nodes[focusIndex] ?? graph.nodes[0]
  const outerNodes = graph.nodes.filter((node) => node.id !== focusNode?.id)
  const center = { x: 380, y: 260 }

  const positions = new Map<string, GraphPoint>()

  if (focusNode) {
    positions.set(focusNode.id, center)
  }

  outerNodes.forEach((node, index) => {
    const angle = (index / Math.max(outerNodes.length, 1)) * Math.PI * 2
    const radius = 180 + (index % 2) * 38
    positions.set(node.id, {
      x: center.x + Math.cos(angle) * radius,
      y: center.y + Math.sin(angle) * radius,
    })
  })

  return positions
}

export const computeViewBox = (points: GraphPoint[]): string => {
  if (points.length === 0) {
    return '0 0 760 520'
  }

  const xs = points.map((point) => point.x)
  const ys = points.map((point) => point.y)
  const minX = Math.min(...xs) - 100
  const maxX = Math.max(...xs) + 100
  const minY = Math.min(...ys) - 90
  const maxY = Math.max(...ys) + 90

  return `${minX} ${minY} ${Math.max(300, maxX - minX)} ${Math.max(240, maxY - minY)}`
}
