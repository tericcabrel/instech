import type { ToolGraphResponseOutput } from '@/api/generated/model/ToolGraphResponse.zod'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'

import { kindLabel } from '../shared/tool-graph-layout'

type ToolGraphSidePanelProps = {
  graph: ToolGraphResponseOutput
  selectedNodeId: string
}

export function ToolGraphSidePanel({ graph, selectedNodeId }: ToolGraphSidePanelProps) {
  const node = graph.nodes.find((item) => item.id === selectedNodeId) ?? graph.nodes[0]

  if (!node) {
    return null
  }

  const linkedEdges = graph.links.filter(
    (link) => link.source === node.id || link.target === node.id,
  )

  return (
    <aside className="feature-panel feature-side-panel">
      <p className="text-muted-foreground mb-1 text-xs uppercase">Node details</p>
      <h2 className="mb-2 text-lg font-semibold">{node.name}</h2>
      <div className="mb-3 flex flex-wrap gap-1.5">
        <Badge variant="secondary">{node.slug}</Badge>
        <Badge>{node.category}</Badge>
        <Badge variant="outline">{node.subType}</Badge>
      </div>

      <dl className="grid grid-cols-2 gap-2 text-sm">
        <dt className="text-muted-foreground">Status</dt>
        <dd className="font-medium">{node.devStatus}</dd>
        <dt className="text-muted-foreground">Release year</dt>
        <dd className="font-medium">{node.releaseYear}</dd>
        <dt className="text-muted-foreground">Degree</dt>
        <dd className="font-medium">{node.degree}</dd>
      </dl>

      <Separator className="my-3" />

      <div className="space-y-2">
        <h3 className="text-xs font-semibold tracking-wide uppercase">Connections</h3>
        {linkedEdges.length === 0 ? (
          <p className="text-muted-foreground text-sm">No relationships for this node.</p>
        ) : (
          <ul className="space-y-2">
            {linkedEdges.map((edge) => {
              const counterpartId = edge.source === node.id ? edge.target : edge.source
              const counterpart = graph.nodes.find((item) => item.id === counterpartId)

              return (
                <li key={edge.id} className="feature-side-connection">
                  <p className="text-sm font-medium">{counterpart?.name ?? counterpartId}</p>
                  <p className="text-muted-foreground text-xs">{kindLabel(edge.kind)}</p>
                  {edge.reason ? <p className="text-xs">{edge.reason}</p> : null}
                </li>
              )
            })}
          </ul>
        )}
      </div>
    </aside>
  )
}
