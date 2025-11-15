import React, { useState } from 'react';
import axios from 'axios';

const UploadPosterForm = ({ animeId, onUpload, onUploadSuccess }) => {
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!file) {
      alert('Please select a file to upload');
      return;
    }

    const formData = new FormData();
    formData.append('poster', file);

    try {
      await axios.post(`http://localhost:8080/api/admin/anime/${animeId}/upload-poster`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      onUpload();
      if (onUploadSuccess) {
        onUploadSuccess();
      }
    } catch (error) {
      console.error('Error uploading poster:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Poster File</label>
        <input type="file" onChange={handleFileChange} />
      </div>
      <button type="submit">Upload</button>
    </form>
  );
};

export default UploadPosterForm;