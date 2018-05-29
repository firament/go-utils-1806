#! /bin/bash
# 2018-05-18 07:44:46 +0530
################################################################################
# v 0.1                                                                        #
#                                                                              #
#    Start godoc server pointing to './src'                                    #
#                                                                              #
# Usage: execute 'bash go-doc' from terminal                                   #
################################################################################
echo "Script to launch go-doc";

# Flag to show only application documentation, or all availaible documentation
#------------------------------------------------------------------------------#
SHOW_APP_DOCS_ONLY=0;
# 0 = Project Only,
# anything else, all availaible docs
#------------------------------------------------------------------------------#

# TODO: Use command arg for setting switch

DOC_PORT=":9048";	# Change this for individual projects
DOC_URL="http://localhost${DOC_PORT}/pkg/?m=all";

# Default Root of the GO project
SOURCE_ROOT="${PWD}/src";
if [[ ${1} && -d ${1} ]]; then
	SOURCE_ROOT=${1};
fi;
# Check to verify ${SOURCE_ROOT} existance
if [[ ! -d ${SOURCE_ROOT} ]]; then
	echo "ERROR: Directory not found!";
	echo "ERROR: ${SOURCE_ROOT}";
	exit 404;
fi;

# Check if server is already running
if [[ ! -z $(pgrep -a "godoc" | grep "${DOC_PORT}") ]]; then
	echo "Port ${DOC_PORT} already in use";
	pgrep -a "godoc" | grep "${DOC_PORT}";
	# Kill only the one on current port
	for gdInst in $(pgrep -a "godoc" | grep "${DOC_PORT}" | cut -d " " -f1); do
		echo "Killing process ${gdInst} to reuse port. ";
		kill -9 ${gdInst}
	done;
fi;

# Apply Scope flag
[[ ! -z ${SHOW_APP_DOCS_ONLY} && ${SHOW_APP_DOCS_ONLY} -eq 0 ]] && GOPATH=${PWD};

# Start godoc server
echo "Starting on port ${DOC_PORT}";
godoc -http=${DOC_PORT} -index -goroot=${SOURCE_ROOT} &

# Open in browser, change if default browser is not desired.
x-www-browser "${DOC_URL}" &
#