import React, { useState } from "react";
import { Box, Typography, Link } from "@mui/material";

export const PreviewTable = ({ images = [], files = [] }) => {
  const [selectedImage, setSelectedImage] = useState(images[0] || null);

  return (
    <Box sx={{ display: "flex", flexDirection: "column", gap: 2, p: 2 }}>
      <Box
        sx={{
          width: "100%",
          height: 300,
          border: "1px solid #ccc",
          borderRadius: 2,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          overflow: "hidden",
        }}
      >
        {selectedImage ? (
          <img
            src={selectedImage}
            alt="preview"
            style={{ maxWidth: "100%", maxHeight: "100%", objectFit: "contain" }}
          />
        ) : (
          <Typography variant="body2" color="text.secondary">
            Нет изображения
          </Typography>
        )}
      </Box>

      <Box sx={{ display: "flex", gap: 1, overflowX: "auto" }}>
        {images.map((img, idx) => (
          <Box
            key={idx}
            sx={{
              border: img === selectedImage ? "2px solid blue" : "1px solid #ccc",
              borderRadius: 1,
              cursor: "pointer",
              flex: "0 0 auto",
              width: 80,
              height: 80,
              overflow: "hidden",
            }}
            onClick={() => setSelectedImage(img)}
          >
            <img
              src={img}
              alt={`thumb-${idx}`}
              style={{ width: "100%", height: "100%", objectFit: "cover" }}
            />
          </Box>
        ))}
      </Box>

      <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
        {files.map((file, i) => (
          <Link key={i} href={file} target="_blank" rel="noopener">
            файл {i + 1}
          </Link>
        ))}
      </Box>
    </Box>
  );
}
