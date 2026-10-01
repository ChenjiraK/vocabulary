"use client";

import { useQuery } from "@tanstack/react-query";
import {
  Button,
  Chip,
  LinearProgress,
  List,
  ListItem,
  ListItemText,
  Paper,
  Typography,
} from "@mui/material";
import { fetchVocabulary } from "../data/vocabulary";

export function VocabularyDemo() {
  const { data, isFetching, refetch } = useQuery({
    queryKey: ["vocabulary-preview"],
    queryFn: fetchVocabulary,
  });

  return (
    <Paper elevation={0} className="border border-slate-200 p-6">
      <div className="flex flex-col gap-6">
        <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
          <div>
            <Typography variant="h5" className="font-bold">
              Vocabulary
            </Typography>
            <Typography color="text.secondary">
              Next.js, Material UI, Tailwind CSS, and TanStack Query are ready.
            </Typography>
          </div>
          <Button variant="contained" onClick={() => void refetch()}>
            Refresh
          </Button>
        </div>

        {isFetching ? <LinearProgress /> : null}

        <List disablePadding>
          {data?.slice(0, 3).map((item) => (
            <ListItem
              key={item.id}
              className="mb-3 rounded-lg border border-slate-200 bg-white"
              secondaryAction={<Chip label={item.level} size="small" />}
            >
              <ListItemText
                primary={
                  <Typography className="font-bold">{item.word}</Typography>
                }
                secondary={item.meaning}
              />
            </ListItem>
          ))}
        </List>
      </div>
    </Paper>
  );
}
