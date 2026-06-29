import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';

type ToolDetailOverviewCardProps = {
  categoryLabel: string;
  details?: string;
  devStatusLabel: string;
  github?: string;
  name: string;
  prolang?: string;
  releaseYear: number;
  slug: string;
  subTypeLabel: string;
  tags: string[];
  useCases: string[];
  website?: string;
};

const renderList = (values: string[]): string => (values.length > 0 ? values.join(', ') : 'None listed');

export const ToolDetailOverviewCard = ({
  categoryLabel,
  details,
  devStatusLabel,
  github,
  name,
  prolang,
  releaseYear,
  slug,
  subTypeLabel,
  tags,
  useCases,
  website,
}: ToolDetailOverviewCardProps) => (
  <Card className="feature-panel p-0">
    <CardHeader className="space-y-2">
      <CardTitle>{name}</CardTitle>
      <p className="text-muted-foreground text-xs">{slug}</p>
      <div className="flex flex-wrap gap-1.5">
        <Badge variant="secondary">{categoryLabel}</Badge>
        <Badge variant="outline">{subTypeLabel}</Badge>
        <Badge variant="outline">{devStatusLabel}</Badge>
      </div>
    </CardHeader>

    <CardContent className="space-y-3 text-xs">
      <div className="grid gap-2 sm:grid-cols-2">
        <div>
          <p className="text-muted-foreground">Release year</p>
          <p>{releaseYear}</p>
        </div>
        <div>
          <p className="text-muted-foreground">Primary language</p>
          <p>{prolang ?? 'Not specified'}</p>
        </div>
      </div>

      <Separator />

      <div className="space-y-1">
        <p className="text-muted-foreground">Details</p>
        <p>{details?.trim() ? details : 'No additional details provided.'}</p>
      </div>

      <Separator />

      <div className="space-y-1">
        <p className="text-muted-foreground">Use cases</p>
        <p>{renderList(useCases)}</p>
      </div>

      <div className="space-y-1">
        <p className="text-muted-foreground">Tags</p>
        <p>{renderList(tags)}</p>
      </div>

      <Separator />

      <div className="grid gap-2 sm:grid-cols-2">
        <div className="space-y-1">
          <p className="text-muted-foreground">Website</p>
          {website ? (
            <a className="underline underline-offset-2" href={website} rel="noreferrer" target="_blank">
              {website}
            </a>
          ) : (
            <p>Not provided</p>
          )}
        </div>
        <div className="space-y-1">
          <p className="text-muted-foreground">GitHub</p>
          {github ? (
            <a className="underline underline-offset-2" href={github} rel="noreferrer" target="_blank">
              {github}
            </a>
          ) : (
            <p>Not provided</p>
          )}
        </div>
      </div>
    </CardContent>
  </Card>
);
