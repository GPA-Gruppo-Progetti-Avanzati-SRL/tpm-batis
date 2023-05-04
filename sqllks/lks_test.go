package sqllks_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqllks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnect(t *testing.T) {

	lks, err := sqllks.GetLinkedService("default")
	require.NoError(t, err)

	_, err = lks.DB()
	require.NoError(t, err)

}

const (
	PQSEQNAME   = "gect_cpx_seq"
	PQSEQPREFIX = "PFX"
)

func TestSequence(t *testing.T) {

	lks, err := sqllks.GetLinkedService("default")
	require.NoError(t, err)

	seq, err := lks.SequenceNextVal(PQSEQNAME, PQSEQPREFIX)
	if err != nil {
		require.NoError(t, err)
	}

	t.Logf("sequence %s has %s value", PQSEQNAME, seq)

}
