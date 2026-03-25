import { useEffect, useEffectEvent, useRef, useState } from 'react';
import { makeAuthenticatedRequest } from '../../auth';

const debounceMilliseconds = 500;
const savedMessageMilliseconds = 2000;
const defaultPonziFactor = 0;

export const minPonziFactor = -0.05;
export const maxPonziFactor = 0.15;

export type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';

interface UsePonziFactorAutosaveProps {
  bankId: string;
  initialPonziFactor: number;
  onSaved: (ponziFactor: number) => void;
}

const normalizePonziFactor = (value: number, fallbackValue: number, source: string) => {
  if (!Number.isFinite(value)) {
    console.error(`Ponzi factor was invalid in ${source}.`, value);
    return fallbackValue;
  }

  return Math.max(minPonziFactor, Math.min(maxPonziFactor, value));
};

export function usePonziFactorAutosave({
  bankId,
  initialPonziFactor,
  onSaved,
}: UsePonziFactorAutosaveProps) {
  const normalizedInitialPonziFactor = normalizePonziFactor(
    initialPonziFactor,
    defaultPonziFactor,
    'props',
  );

  const [draftPonziFactor, setDraftPonziFactor] = useState(normalizedInitialPonziFactor);
  const [savedPonziFactor, setSavedPonziFactor] = useState(normalizedInitialPonziFactor);
  const [saveStatus, setSaveStatus] = useState<SaveStatus>('idle');

  const debounceTimeoutIdRef = useRef<number | null>(null);
  const savedMessageTimeoutIdRef = useRef<number | null>(null);
  const activeRequestIdRef = useRef(0);
  const isSaveInProgressRef = useRef(false);
  const previousBankIdRef = useRef(bankId);
  const draftPonziFactorRef = useRef(normalizedInitialPonziFactor);
  const savedPonziFactorRef = useRef(normalizedInitialPonziFactor);

  function clearDebounceTimeout() {
    if (debounceTimeoutIdRef.current !== null) {
      window.clearTimeout(debounceTimeoutIdRef.current);
      debounceTimeoutIdRef.current = null;
    }
  }

  function clearSavedMessageTimeout() {
    if (savedMessageTimeoutIdRef.current !== null) {
      window.clearTimeout(savedMessageTimeoutIdRef.current);
      savedMessageTimeoutIdRef.current = null;
    }
  }

  function setCurrentDraftPonziFactor(nextPonziFactor: number) {
    draftPonziFactorRef.current = nextPonziFactor;
    setDraftPonziFactor(nextPonziFactor);
  }

  function setCurrentSavedPonziFactor(nextPonziFactor: number) {
    savedPonziFactorRef.current = nextPonziFactor;
    setSavedPonziFactor(nextPonziFactor);
  }

  const showSavedStatus = useEffectEvent(() => {
    clearSavedMessageTimeout();
    setSaveStatus('saved');

    savedMessageTimeoutIdRef.current = window.setTimeout(() => {
      savedMessageTimeoutIdRef.current = null;

      setSaveStatus((currentStatus) => {
        if (currentStatus !== 'saved') {
          return currentStatus;
        }

        return 'idle';
      });
    }, savedMessageMilliseconds);
  });

  const performSave = useEffectEvent(async (valueToSave: number, bankIdToSave: string) => {
    const response = await makeAuthenticatedRequest('/api/bank/ponziFactor', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ bankId: bankIdToSave, ponziFactor: valueToSave }),
    });

    if (!response.ok) {
      throw new Error('Save failed');
    }
  });

  const savePonziFactor = useEffectEvent(async (valueToSave: number) => {
    if (isSaveInProgressRef.current) {
      return;
    }

    clearSavedMessageTimeout();
    isSaveInProgressRef.current = true;
    setSaveStatus('saving');
    const requestId = activeRequestIdRef.current + 1;
    activeRequestIdRef.current = requestId;

    try {
      await performSave(valueToSave, bankId);

      if (requestId !== activeRequestIdRef.current) {
        return;
      }

      isSaveInProgressRef.current = false;
      setCurrentSavedPonziFactor(valueToSave);
      onSaved(valueToSave);

      if (draftPonziFactorRef.current !== valueToSave) {
        void savePonziFactor(draftPonziFactorRef.current);
        return;
      }

      showSavedStatus();
    } catch {
      if (requestId !== activeRequestIdRef.current) {
        return;
      }

      isSaveInProgressRef.current = false;
      setSaveStatus('error');
    }
  });

  const requestSave = useEffectEvent((valueToSave: number) => {
    clearSavedMessageTimeout();

    if (valueToSave === savedPonziFactorRef.current && !isSaveInProgressRef.current) {
      setSaveStatus('idle');
      return;
    }

    setSaveStatus('saving');

    if (isSaveInProgressRef.current) {
      return;
    }

    void savePonziFactor(valueToSave);
  });

  const setPonziFactorDraft = useEffectEvent((nextValue: number) => {
    const normalizedNextPonziFactor = normalizePonziFactor(
      nextValue,
      draftPonziFactorRef.current,
      'range input',
    );

    clearDebounceTimeout();
    setCurrentDraftPonziFactor(normalizedNextPonziFactor);

    if (normalizedNextPonziFactor === savedPonziFactorRef.current) {
      if (!isSaveInProgressRef.current) {
        clearSavedMessageTimeout();
        setSaveStatus('idle');
        return;
      }

      setSaveStatus('saving');
      return;
    }

    clearSavedMessageTimeout();
    setSaveStatus('saving');

    debounceTimeoutIdRef.current = window.setTimeout(() => {
      debounceTimeoutIdRef.current = null;
      requestSave(normalizedNextPonziFactor);
    }, debounceMilliseconds);
  });

  const flushPendingSave = useEffectEvent(() => {
    if (debounceTimeoutIdRef.current !== null) {
      clearDebounceTimeout();
      requestSave(draftPonziFactorRef.current);
      return;
    }

    if (draftPonziFactorRef.current !== savedPonziFactorRef.current || isSaveInProgressRef.current) {
      requestSave(draftPonziFactorRef.current);
    }
  });

  useEffect(() => {
    if (previousBankIdRef.current === bankId) {
      return;
    }

    previousBankIdRef.current = bankId;
    const normalizedPropValue = normalizePonziFactor(initialPonziFactor, defaultPonziFactor, 'props');
    clearSavedMessageTimeout();
    clearDebounceTimeout();
    activeRequestIdRef.current += 1;
    isSaveInProgressRef.current = false;
    setCurrentDraftPonziFactor(normalizedPropValue);
    setCurrentSavedPonziFactor(normalizedPropValue);
    setSaveStatus('idle');
  }, [bankId, initialPonziFactor]);

  useEffect(() => {
    if (debounceTimeoutIdRef.current !== null) {
      return;
    }

    if (isSaveInProgressRef.current) {
      return;
    }

    if (draftPonziFactor !== savedPonziFactor) {
      return;
    }

    const normalizedPropValue = normalizePonziFactor(
      initialPonziFactor,
      savedPonziFactor,
      'props',
    );

    if (normalizedPropValue === draftPonziFactor && normalizedPropValue === savedPonziFactor) {
      return;
    }

    clearSavedMessageTimeout();
    isSaveInProgressRef.current = false;
    setCurrentDraftPonziFactor(normalizedPropValue);
    setCurrentSavedPonziFactor(normalizedPropValue);
    setSaveStatus('idle');
  }, [bankId, draftPonziFactor, initialPonziFactor, savedPonziFactor]);

  useEffect(() => {
    return () => {
      activeRequestIdRef.current += 1;
      isSaveInProgressRef.current = false;
      clearDebounceTimeout();
      clearSavedMessageTimeout();
    };
  }, []);

  return {
    draftPonziFactor,
    saveStatus,
    setPonziFactorDraft,
    flushPendingSave,
  };
}