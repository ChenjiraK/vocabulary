"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Chip,
  LinearProgress,
  MenuItem,
  Paper,
  Select,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import { fetchVocabulary, type VocabularyLevel } from "../data/vocabulary";

const levels: Array<"All" | VocabularyLevel> = ["All", "A1", "A2", "B1", "B2"];

export function VocabularyTable() {
  const [search, setSearch] = useState("");
  const [level, setLevel] = useState<"All" | VocabularyLevel>("All");

  const { data = [], isFetching } = useQuery({
    queryKey: ["vocabulary-table"],
    queryFn: fetchVocabulary,
  });

  const filteredVocabulary = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase();

    return data.filter((item) => {
      const matchesLevel = level === "All" || item.level === level;
      const matchesSearch =
        normalizedSearch.length === 0 ||
        item.word.toLowerCase().includes(normalizedSearch) ||
        item.meaning.toLowerCase().includes(normalizedSearch) ||
        item.partOfSpeech.toLowerCase().includes(normalizedSearch);

      return matchesLevel && matchesSearch;
    });
  }, [data, level, search]);

  return (
    <Paper elevation={0} className="overflow-hidden border border-slate-200">
      <div className="flex flex-col gap-4 border-b border-slate-200 p-5 md:flex-row md:items-center md:justify-between">
        <div>
          <Typography variant="h6" className="font-bold">
            Vocabulary List
          </Typography>
          <Typography color="text.secondary">
            {filteredVocabulary.length} words found
          </Typography>
        </div>

        <div className="grid gap-3 sm:grid-cols-[minmax(220px,1fr)_130px]">
          <TextField
            label="Search"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            size="small"
          />
          <Select
            value={level}
            onChange={(event) =>
              setLevel(event.target.value as "All" | VocabularyLevel)
            }
            size="small"
            displayEmpty
          >
            {levels.map((item) => (
              <MenuItem key={item} value={item}>
                {item}
              </MenuItem>
            ))}
          </Select>
        </div>
      </div>

      {isFetching ? <LinearProgress /> : null}

      <TableContainer>
        <Table>
          <TableHead>
            <TableRow className="bg-slate-50">
              <TableCell>Word</TableCell>
              <TableCell>Meaning</TableCell>
              <TableCell>Part of Speech</TableCell>
              <TableCell>Level</TableCell>
              <TableCell>Example</TableCell>
              <TableCell>Last Reviewed</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredVocabulary.map((item) => (
              <TableRow key={item.id} hover>
                <TableCell>
                  <Typography className="font-bold">{item.word}</Typography>
                </TableCell>
                <TableCell>{item.meaning}</TableCell>
                <TableCell>{item.partOfSpeech}</TableCell>
                <TableCell>
                  <Chip label={item.level} size="small" color="primary" />
                </TableCell>
                <TableCell className="max-w-xs">{item.example}</TableCell>
                <TableCell>{item.lastReviewed}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {filteredVocabulary.length === 0 ? (
        <div className="p-8 text-center text-slate-500">
          No vocabulary found.
        </div>
      ) : null}
    </Paper>
  );
}
