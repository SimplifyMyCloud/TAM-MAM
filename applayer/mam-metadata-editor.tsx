import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus, X, ChevronDown, ChevronUp } from 'lucide-react';

const MetadataEditor = ({ asset, onSave }) => {
  const [metadata, setMetadata] = useState({});
  const [expandedSections, setExpandedSections] = useState({
    technical: true,
    descriptive: true,
    rights: true,
    custom: true
  });

  useEffect(() => {
    if (asset?.metadata) {
      setMetadata(asset.metadata);
    }
  }, [asset]);

  const handleMetadataChange = (section, field, value) => {
    setMetadata(prev => ({
      ...prev,
      [section]: {
        ...prev[section],
        [field]: value
      }
    }));
  };

  const handleAddCustomField = (section) => {
    const fieldName = prompt("Enter field name:");
    if (fieldName) {
      setMetadata(prev => ({
        ...prev,
        [section]: {
          ...prev[section],
          [fieldName]: ""
        }
      }));
    }
  };

  const handleRemoveField = (section, field) => {
    setMetadata(prev => {
      const sectionData = { ...prev[section] };
      delete sectionData[field];
      return {
        ...prev,
        [section]: sectionData
      };
    });
  };

  const toggleSection = (section) => {
    setExpandedSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }));
  };

  return (
    <div className="space-y-4">
      {/* Technical Metadata */}
      <Card>
        <CardHeader className="cursor-pointer" onClick={() => toggleSection('technical')}>
          <div className="flex justify-between items-center">
            <CardTitle className="text-lg">Technical Metadata</CardTitle>
            {expandedSections.technical ? <ChevronUp /> : <ChevronDown />}
          </div>
        </CardHeader>
        {expandedSections.technical && (
          <CardContent className="space-y-4">
            <MetadataField
              label="Format"
              value={metadata.technical?.format || ""}
              onChange={(value) => handleMetadataChange('technical', 'format', value)}
            />
            <MetadataField
              label="Duration"
              value={metadata.technical?.duration || ""}
              onChange={(value) => handleMetadataChange('technical', 'duration', value)}
            />
            <MetadataField
              label="Resolution"
              value={metadata.technical?.resolution || ""}
              onChange={(value) => handleMetadataChange('technical', 'resolution', value)}
            />
            <MetadataField
              label="Codec"
              value={metadata.technical?.codec || ""}
              onChange={(value) => handleMetadataChange('technical', 'codec', value)}
            />
            <MetadataField
              label="Bitrate"
              value={metadata.technical?.bitrate || ""}
              onChange={(value) => handleMetadataChange('technical', 'bitrate', value)}
            />
          </CardContent>
        )}
      </Card>

      {/* Descriptive Metadata */}
      <Card>
        <CardHeader className="cursor-pointer" onClick={() => toggleSection('descriptive')}>
          <div className="flex justify-between items-center">
            <CardTitle className="text-lg">Descriptive Metadata</CardTitle>
            {expandedSections.descriptive ? <ChevronUp /> : <ChevronDown />}
          </div>
        </CardHeader>
        {expandedSections.descriptive && (
          <CardContent className="space-y-4">
            <MetadataField
              label="Title"
              value={metadata.descriptive?.title || ""}
              onChange={(value) => handleMetadataChange('descriptive', 'title', value)}
            />
            <div className="space-y-2">
              <label className="block text-sm font-medium">Description</label>
              <textarea
                className="w-full min-h-[100px] p-2 border rounded-md"
                value={metadata.descriptive?.description || ""}
                onChange={(e) => handleMetadataChange('descriptive', 'description', e.target.value)}
              />
            </div>
            <MetadataField
              label="Keywords"
              value={metadata.descriptive?.keywords || ""}
              onChange={(value) => handleMetadataChange('descriptive', 'keywords', value)}
              placeholder="Comma-separated keywords"
            />
            <MetadataField
              label="Language"
              value={metadata.descriptive?.language || ""}
              onChange={(value) => handleMetadataChange('descriptive', 'language', value)}
            />
          </CardContent>
        )}
      </Card>

      {/* Rights Metadata */}
      <Card>
        <CardHeader className="cursor-pointer" onClick={() => toggleSection('rights')}>
          <div className="flex justify-between items-center">
            <CardTitle className="text-lg">Rights Management</CardTitle>
            {expandedSections.rights ? <ChevronUp /> : <ChevronDown />}
          </div>
        </CardHeader>
        {expandedSections.rights && (
          <CardContent className="space-y-4">
            <MetadataField
              label="Copyright"
              value={metadata.rights?.copyright || ""}
              onChange={(value) => handleMetadataChange('rights', 'copyright', value)}
            />
            <MetadataField
              label="License"
              value={metadata.rights?.license || ""}
              onChange={(value) => handleMetadataChange('rights', 'license', value)}
            />
            <MetadataField
              label="Usage Restrictions"
              value={metadata.rights?.restrictions || ""}
              onChange={(value) => handleMetadataChange('rights', 'restrictions', value)}
            />
            <div className="space-y-2">
              <label className="block text-sm font-medium">Expiry Date</label>
              <Input
                type="date"
                value={metadata.rights?.expiryDate || ""}
                onChange={(e) => handleMetadataChange('rights', 'expiryDate', e.target.value)}
              />
            </div>
          </CardContent>
        )}
      </Card>

      {/* Custom Metadata */}
      <Card>
        <CardHeader className="cursor-pointer" onClick={() => toggleSection('custom')}>
          <div className="flex justify-between items-center">
            <CardTitle className="text-lg">Custom Metadata</CardTitle>
            {expandedSections.custom ? <ChevronUp /> : <ChevronDown />}
          </div>
        </CardHeader>
        {expandedSections.custom && (
          <CardContent className="space-y-4">
            {Object.entries(metadata.custom || {}).map(([field, value]) => (
              <div key={field} className="flex items-center gap-2">
                <MetadataField
                  label={field}
                  value={value}
                  onChange={(newValue) => handleMetadataChange('custom', field, newValue)}
                />
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleRemoveField('custom', field)}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            ))}
            <Button
              variant="outline"
              onClick={() => handleAddCustomField('custom')}
              className="w-full"
            >
              <Plus className="h-4 w-4 mr-2" /> Add Custom Field
            </Button>
          </CardContent>
        )}
      </Card>

      {/* Save Button */}
      <div className="flex justify-end space-x-2">
        <Button
          variant="outline"
          onClick={() => setMetadata(asset?.metadata || {})}
        >
          Reset
        </Button>
        <Button
          onClick={() => onSave(metadata)}
        >
          Save Metadata
        </Button>
      </div>
    </div>
  );
};

// Reusable Metadata Field Component
const MetadataField = ({ label, value, onChange, placeholder }) => {
  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium">{label}</label>
      <Input
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
      />
    </div>
  );
};

export default MetadataEditor;
