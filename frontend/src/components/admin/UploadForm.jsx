
import React, { useState } from 'react';
import axios from 'axios';

const UploadForm = ({ episodeId, onUpload }) => {
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
    formData.append('video', file);

    try {
      await axios.post(`http://localhost:8080/api/admin/episode/${episodeId}/upload-raw-video`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      onUpload();
    } catch (error) {
      console.error('Error uploading video:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Video File</label>
        <input type="file" onChange={handleFileChange} />
      </div>
      <button type="submit">Upload</button>
    </form>
  );
};

export default UploadForm;
