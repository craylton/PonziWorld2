import { useId } from 'react';
import {
  maxPonziFactor,
  minPonziFactor,
  type SaveStatus,
  usePonziFactorAutosave,
} from './usePonziFactorAutosave';

interface PonziFactorSliderProps {
  bankId: string;
  ponziFactor: number;
  onPonziFactorSaved: (ponziFactor: number) => void;
}

const statusText: Record<SaveStatus, string> = {
  idle: '',
  saving: 'Saving...',
  saved: 'Saved',
  error: 'Error',
};

const formatPercent = (value: number) => {
  const percent = value * 100;
  const rounded = Math.round(percent * 10) / 10;

  if (Number.isInteger(rounded)) {
    return `${rounded.toFixed(0)}%`;
  }

  return `${rounded.toFixed(1)}%`;
};

export default function PonziFactorSlider({
  bankId,
  ponziFactor,
  onPonziFactorSaved,
}: PonziFactorSliderProps) {
  const ponziFactorInputId = useId();
  const {
    draftPonziFactor,
    saveStatus,
    setPonziFactorDraft,
    flushPendingSave,
  } = usePonziFactorAutosave({
    bankId,
    initialPonziFactor: ponziFactor,
    onSaved: onPonziFactorSaved,
  });

  const statusModifier = saveStatus !== 'idle' ? ` ponzi-value__status--${saveStatus}` : '';

  return (
    <div className="ponzi-value">
      <label className="ponzi-value__label" htmlFor={ponziFactorInputId}>
        Ponzi factor:{' '}
        <span className="ponzi-value__percent">
          {formatPercent(draftPonziFactor)}
        </span>
      </label>

      <div className="ponzi-value__controls">
        <input
          id={ponziFactorInputId}
          type="range"
          min={minPonziFactor}
          max={maxPonziFactor}
          step={0.001}
          value={draftPonziFactor}
          onChange={(event) => setPonziFactorDraft(Number(event.target.value))}
          onBlur={flushPendingSave}
        />
        <span className={`ponzi-value__status${statusModifier}`} aria-live="polite">
          {statusText[saveStatus]}
        </span>
      </div>
    </div>
  );
}
