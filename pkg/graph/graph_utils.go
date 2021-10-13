// SPDX-FileCopyrightText: 2021 k8sviz authors
// SPDX-License-Identifier: Apache-2.0

package graph

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// imagePath returns the path to the image file
// path is {dir}/icons/{resource}-128.png
// ex) /icons/pod-128.png
func (g *Graph) imagePath(kind string) string {
	return filepath.Join(g.dir, "icons", kind+imageSuffix)
}

func (g *Graph) imagePathRelative(kind string) string {
	return filepath.Join("icons", kind+imageSuffix)
}

func (g *Graph) imagePathBase64(kind string) string {
	path := filepath.Join(g.dir, "icons", kind+imageSuffix)

	if b64string, ok := g.base64Icons[path]; ok {
		return "data:image/png;charset=utf-8;base64," + b64string
	}

	// Open file on disk.
	f, _ := os.Open(path)

	// Read entire file into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return "data:image/png;charset=utf-8;base64," + encoded
}

// clusterLabel returns the resource label for namespace
// ex)
//   <<TABLE BORDER="0"><TR><TD><IMG SRC="/icons/ns-128.png" /></TD></TR><TR><TD>my-namespace</TD></TR></TABLE>>
func (g *Graph) clusterLabel() string {
	return g.resourceLabel("ns", g.res.Namespace)
}

// resourceLabel returns the resource label for a resource
// ex)
//   <<TABLE BORDER="0"><TR><TD><IMG SRC="/icons/pod-128.png" /></TD></TR><TR><TD>my-pod</TD></TR></TABLE>>
func (g *Graph) resourceLabel(kind, name string) string {
	return fmt.Sprintf("<<TABLE BORDER=\"0\"><TR><TD><IMG SRC=\"%s\" /></TD></TR><TR><TD>%s</TD></TR></TABLE>>", g.imagePathRelative(kind), name)
}

// clusterName returns name of the graphviz cluster
// It is named base on namespace.
// ex) cluster_my_namespace
func (g *Graph) clusterName() string {
	return clusterPrefix + g.escapeName(g.res.Namespace)
}

// escapeName returns the escaped name to be handled with graphviz
// It replaces "." and "-" with "_".
// ex) my_namespace
func (g *Graph) escapeName(name string) string {
	return strings.NewReplacer(".", "_", "-", "_").Replace(name)
}

// resourceName returns the escaped name of the resource
// It espaces the resource name and add resType as a prefix.
// ex) pod_my_pod
func (g *Graph) resourceName(resType, name string) string {
	return resType + "_" + g.escapeName(name)
}

// rankName returns the name of the dummy rank
// ex) rank_1
func (g *Graph) rankName(rank int) string {
	return fmt.Sprintf("%s%d", rankPrefix, rank)
}

// rankDummyNodeName returns the node name of the dummy rank
// ex) 1
func (g *Graph) rankDummyNodeName(rank int) string {
	return fmt.Sprintf("%d", rank)
}
